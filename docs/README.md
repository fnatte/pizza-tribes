# Pizza Tribes Documentation

This document provides implementation details for each feature. More specifically, it shows how Redis is used and **the Redis commands to store and retrieve data.**

## Table of Contents

* [Pizza Tribes Documentation](#pizza-tribes-documentation)
   * [Users — Registration and Authentication](#users--registration-and-authentication)
   * [Game State](#game-state)
   * [Game State Update](#game-state-update)
      * [Figuring out Whom Needs Update](#figuring-out-whom-needs-update)
      * [Updating the Game State](#updating-the-game-state)
      * [Schedule for Next Update](#schedule-for-next-update)
      * [Insert Data Points (Timeseries)](#insert-data-points-timeseries)
      * [Updating the Leaderboard](#updating-the-leaderboard)

## Users &mdash; Registration and Authentication

The users are stored as a hash set in key `user:{user_id}` containing fields:
- `id`
- `username`
- `hashed_password`

The `user_id` can be looked up using the username via `username:{username}`.

**Registration** is done like so:

- Generate unique id (rs/xid)
- redis cmd: `SET username:{username} user_id`
- redis cmd: `HSET user:{user_id} "id" user_id "username" username "hashed_password" hash`

**Authentication** is done like so:

- redis cmd: `GET username:{username}` (get user id)
- redis cmd: `HGETALL user:{user_id}`
- Verify hashed_password
	- If the hashes match, create JWT

## Game State

The user game state is stored as a JSON value (using RedisJSON) in key `user:{user_id}:gamestate` with the following structure:

```json
{
  "resources": {
    "coins": 0,
    "pizzas": 0
  },
  "lots": {},
  "population": {
    "uneducated": 0,
    "chefs": 0,
    "salesmice": 0,
    "guards": 0,
    "thieves": 0
  },
  "timestamp": 1620842714,
  "trainingQueue": [],
  "constructionQueue": [],
  "townX": 50,
  "townY": 50,
  "travelQueue": []
}
```

See [/protos/gamestate.proto](/protos/gamestate.proto) for full definition.

The game data is accessed in different ways depending on the use case. But for a complete retrieval, the following is used:

- redis cmd: `JSON.GET user:{user_id}.gamestate .`

In other cases a path is used to retrieve only a subset of the data:

- redis cmd: `JSON.GET user:{user_id}.gamestate '.lots["5"]'` (retrieve building info at lot 5)
- redis cmd: `JSON.GET user:{user_id}.gamestate .population` (retrieve population data)

## Game State Update

The game state update is what makes the game tick. It is one of the most important processes in the game. The purpose of a game state update is to:
- Extrapolate resources (i.e., increase resources with produced amounts since the last update)
- Complete buildings
- Complete trainings
- Complete travels (i.e., thieves moving between towns)

It will also:
- Insert resource data points for timeseries
- Update leaderboard (because the resources have changed)

> ℹ️ Notice that the game could have been implemented using a Redis stream of game events, and as such, the game state (and snapshots) could have been derived from the stream of events. That idea was not pursued because I estimated that the fastest way to build the game was using a simple game state update loop.

### Figuring out Whom Needs Update

The updater runs in a loop that queries a sorted set named `user_updates`. It retrieves the top record in the sorted set by running:

```
ZRANGE user_updates 0 0 WITHSCORES
{1.6208459243016696e+18 c2e16taink8s73ejr3qg}
```

By utilizing `WITHSCORES` we also retrieve the timestamp of when that user needs a game state. As such, the updater can check if `timestamp < now`, and if so:

1. `ZREM user_updates c2e19af8q04s73f8j8lg`
2. Proceed to update the game state

> Note that this is not entirely safe because if the game state update fails, the user will no longer have a record in the `user_updates` sorted set, and as such, it will not be scheduled for any game state update.
> As a workaround, the game currently ensures that the user is scheduled for game state updates on login.


### Updating the Game State

The update is executed with a check-and-set approach (`WATCH`, `MULTI`, `EXEC`):
1. `WATCH user:{user_id}:gamestate`
1. `JSON.GET user:{user_id}:gamestate`
1. Run game state process to figure out how to transform the game state
1. `MULTI`
1. Run all modifying commands to transformed to the game state calculated in previous step
1. `EXEC`


See this trace (of a simple game state update) for details:

```text
watch user:c2e19af8q04s73f8j8lg:gamestate: OK
JSON.GET user:c2e19af8q04s73f8j8lg:gamestate .: {"resources":{"coins":20,"pizzas":0},"lots":{"1":{"building":3},"2":{"building":0},"9":{"building":1},"10":{"building":2}},"population":{"uneducated":8,"chefs":1,"salesmice":1,"guards":0,"thieves":0},"timestamp":1620845911,"trainingQueue":[],"constructionQueue":[],"townX":51,"townY":58,"travelQueue":[]}
[multi: QUEUED
	JSON.SET user:c2e19af8q04s73f8j8lg:gamestate .timestamp 1620845921: OK
	JSON.SET user:c2e19af8q04s73f8j8lg:gamestate .resources.coins 22: OK
	JSON.SET user:c2e19af8q04s73f8j8lg:gamestate .resources.pizzas 0: OK
	exec: []]
unwatch: OK
```


Extrapolating resources, completing buildings, completing trainings, and complete travels are all implemented using the flow described above. The difference is what modifying commands are queued to manipulate the game state.

### Schedule for Next Update

When the game state has been updated we must also schedule the next one:

1. Determine when the game state needs to be updated
	1. Is a building being completed?
	1. Is a training being completed?
	1. Is a travel being completed?
1. `ZADD user_updates {timestamp_of_next_update} {user_id}`

### Insert Data Points (Timeseries)

The RedisTimeseries module is used to track the changes in user resources. The resources are tracked using the following keys:

- `user:c2e19af8q04s73f8j8lg:ts_coins`
- `user:c2e19af8q04s73f8j8lg:ts_pizzas`

Upon every game state update a new data point is inserted into each key like so:

`TS.ADD user:{user_id}:ts_coins {timestamp_now} {current_amount_of_coins}`
`TS.ADD user:{user_id}:ts_pizzas {timestamp_now} {current_amount_of_pizzas}`

When the user wants to look at their resource history, the following command is used to retrieve the aggregated data points from the last 24 hours:

```
	from := now - 24 hours
	to := now
	timeBucket := 1 hour

	TS.RANGE user:{user_id}:ts_coins {from} {to} AGGREGATION avg {timeBucket}
	TS.RANGE user:{user_id}:ts_pizzas {from} {to} AGGREGATION avg {timeBucket}
```

### Updating the Leaderboard

The game state update will change the number of coins a user has. That is why we need to update the leaderboard.

The leaderboard is a sorted set with the key `leaderboard`. It is updated by running:

```
ZADD leaderboard {current_amount_of_coins} {user_id}
```

When any user wants to take a look at the leaderboard, the data is retrieved like so:

```
ZREVRANGE leaderboard 0 20 WITHSCORES
```

