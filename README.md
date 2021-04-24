# Pizza Mouse

This project is aimed to be registered for [RedisConf 2021
Hackathon](https://hackathons.redislabs.com/hackathons/build-on-redis-hackathon).

## Project Idea

Pizza Mouse is a multiplayer persistent browser-based clicker real-time
strategy game. The gameplay is a combination of a clicker game
(https://en.wikipedia.org/wiki/Incremental_game) and a real-time strategy
game
(https://en.wikipedia.org/wiki/Massively_multiplayer_online_real-time_strategy_game).

As a player, you are in charge of building and expanding your pizzeria by
building and assigning roles to your mice population. Your pizza empire
will earn coins by selling pizzas. The player with the most coins is
proclaimed as the winner.

To keep your profits intact, you must train part of your population as
guards to protect yourself from being robbed by other players. And maybe,
you join the dirty game yourself and send thieves on your competitors!

## Game Elements

### Buildings

- Kitchen
- Shop
- House
- School

### Education / Roles

- Chef
- Salesmouse
- Guard
- Thief

### Resources

- Pizzas
- Coins

### Actions

- Click (tap)
- Construct building
- Upgrade building
- Train
- Expand
- Steal

### Other Features

- Leaderboard
- Reports

## Architecture and Use of Redis


### Overview

```
          +---------+
          | Web App |
          +---------+
               |
               | (Web Socket / HTTPS)
               |
          +---------+
          | Web Api |
          +---------+
               |
               |
               |
          +---------+
    +-----|  Redis  |-----+
    |     +---------+     |
    |                     |
    |                     |
    |                     |
  +--------+         +--------+
  | Worker |         | Worker |
  +--------+         +--------+
```


### Client-server Communication

The project takes an (probably unconventional) approach of
client-server communication, relying heavily on web sockets instead
even for communication that traditionally is fulfilled by HTTP
request/response.

1. The _Web App_ sends commands over the Web socket
2. The _Web API_ enqueue the command on a Redis queue _in_
3. A _Worker_
	a. Pulls the command from the Redis queue _in_
	b. Executes the command
	c. May push a response to another Redis queue _out_
4. The _Web API_
	a. Pulls a response from the Redis queue _out_
	b. Sends the response back to the corresponding Web socket

#### Web Sockets

Web socket communication heavily relies on Redis for pushing and pulling
messages. The API does not do any game logic but simply validates client
messages before pushing them to Redis. This allows for running multiple
game workers to do the "heavy" lifting (not that heavy really :)). The
workers can be horizontally scaled while relying on Redis performance to
push messages through the system. Since Web Sockets are stateful
bidirectional communication over a single TCP connection, it is not easy
to scale the _Web API_ (holding the sockets). This solution attempts to
minimize load on the _Web API_ so that it can focus on shoveling data to
the clients.

### Delayed Tasks

The worker needs to delay some tasks (e.g., finish construction of
building after 5 minutes). This is accomplished by pushing to the sorted
set _delayed_. A worker will pick up the delayed task by polling the
sorted set.

### Client-Server Protocol

#### Client Messages

| Id  |  Type                 |  Payload           |
|-----|-----------------------|--------------------|
| ... |  "TAP"                | amount?            |
| ... |  "CONSTRUCT_BUILDING" | lotId, building    |
| ... |  "UPGRADE_BUILDING"   | lotId              |
| ... |  "TRAIN"              | education, amount  |
| ... |  "EXPAND"             |                    |
| ... |  "STEAL"              | amount, x, y       |

#### Server Messages

|  Type                 |  Payload           |
|-----------------------|--------------------|
|  "STATE_CHANGE"       | ...game_state      |
|  "RESPONSE"           | request_id, result |

## Tech Stack

- The Web App is/will be built with _TypeScript_, _React_, _Zustand_
- The Web API and workers are/will be built with _Go_ and _Redis_


The backend should strive for using Redis using go-redis only,
implementing queues, locks, delayed tasks, etc., without additional
libraries.

