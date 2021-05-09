# Pizza Tribes

Play at: https://pizzatribes.teus.dev

This project is aimed to be registered for [RedisConf 2021
Hackathon](https://hackathons.redislabs.com/hackathons/build-on-redis-hackathon).

## Project Idea

Pizza Tribes is a multiplayer persistent browser-based clicker real-time
strategy game. The gameplay is a combination of a [clicker
game](https://en.wikipedia.org/wiki/Incremental_game) and a [real-time
strategy
game](https://en.wikipedia.org/wiki/Massively_multiplayer_online_real-time_strategy_game).

As a player, you are in charge of building and expanding your town by
building and assigning roles to your mice population. Your pizza empire
will earn coins by selling pizzas. The player with the most coins is
proclaimed as the winner.

To keep your profits intact, you must train part of your population as
guards to protect yourself from being robbed by other players. And maybe,
you join the dirty game yourself and send thieves on your competitors!

## Game Elements

### Buildings

- [x] Kitchen
- [x] Shop
- [x] House
- [x] School

### Education / Roles

- [x] Chef
- [x] Salesmouse
- [ ] Guard (can be trained but does nothing)
- [x] Thief

### Resources

- [x] Pizzas
- [x] Coins

### Actions

- [ ] Click (tap)
- [x] Construct building
- [x] Upgrade building
- [x] Train
- [x] Steal
- [ ] Expand

### Other Features

- [x] Leaderboard
- [x] Reports

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
    +-----|  Redis  |<-----------+
    |     +---------+            |
    |            |          +---------+
    |            |          | Updater |
    |            |          +---------+
  +--------+   +--------+
  | Worker |   | Worker |
  +--------+   +--------+
```

As the diagram above describe, there are three types of backend services:

- _Web Api_ &mdash; serves HTTP requests and holds Web sockets
- _Worker_ &mdash; processes client messages pulled from Redis
- _Updater_ &mdash; processes delayed game state updates

In addition, the diagram show:

- The _Web App_ &mdash; the game client
- **Redis** &mdash; used as swiss army knife for data communication and persistence


### Client-server Communication

The project takes an (probably unconventional) approach of
client-server communication, relying heavily on web sockets
even for communication that traditionally is fulfilled by HTTP
request/response. However, there are a few traditional REST-like API
endpoint as well.

A typical Web socket flow goes as follows:

1. The _Web App_ sends commands over the Web socket
2. The _Web API_ enqueue the command on a Redis queue _wsin_ (`RPUSH`)
3. A _Worker_
	1. Pulls the command from the Redis queue _wsin_(`BLPop`)
	1. Executes the command
	1. May push a response to another Redis queue _wsout_ (`RPUSH`)
4. The _Web API_
	1. Pulls a response from the Redis queue _wsout_ (`BLPOP`)
	1. Sends the response back to the corresponding Web socket

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

Note that the messages are not sent between the API and workers using
pub/sub but instead using Redis lists (RPush and BLPop). I am not sure if
this is used as a good idea or not - but I think one upside is that the workers or
API can be restarted without losing messages (pub/sub is fire-and-forget).

### Updater (Delayed Tasks)

The worker needs to delay some tasks (e.g., finish construction of
building after 5 minutes). This is accomplished by updating to the sorted
set _user_updates_. The *updater* pulls the top record of the sorted set
(sorted by time), and if the time has passed it removes the record from
the set, and then updates the game state of that user (to finish the
construction).

A (simplified) typical flow is as follows:

1. The *Web App* send command to start construction of a building
1. A *worker* processes the command
	1. Validates the command
	1. Updates the user game state: `JSON.ARRAPPEND user:$user_id:gamestate .constructionQueue $constructionItem`
	1. Finds the next time the user game state needs to be updated (e.g. when the construction is completed)
	1. Set next update time: `ZADD user_updates $timestamp $user_id`
1. A *updater* updates the user game state at the next update time
	1. Runs `ZRANGE user_updates 0 0 WITHSCORES` to fetch the next user that needs update (and at what time)
	1. If the score (timestamp) has been passed
		1. Remove the next update time: `ZREM user_updates $user_id`
		1. Perform game state update
		1. Find next time the user game state needs to be updated again
		1. Set next update time: `ZADD user_updates $timestamp $user_id`

[![](https://mermaid.ink/img/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG4gICAgcGFydGljaXBhbnQgdXBkYXRlclxuICAgIHBhcnRpY2lwYW50IHVzZXJfdXBkYXRlc1xuICAgIHBhcnRpY2lwYW50IGdhbWVfc3RhdGVcblxuICAgIGxvb3BcbiAgICB1cGRhdGVyIC0-PiB1c2VyX3VwZGF0ZXM6IFpSQU5HRSB1c2VyX3VwZGF0ZXMgMCAwIFdJVEhTQ09SRVNcbiAgICB1c2VyX3VwZGF0ZXMgLT4-IHVwZGF0ZXI6ICgkdGltZXN0YW1wLCAkdXNlcl9pZClcbiAgICBhbHQgdGltZXN0YW1wIDwgbm93XG4gICAgICB1cGRhdGVyIC0-PiB1c2VyX3VwZGF0ZXM6IFpSRU0gdXNlcl91cGRhdGVzICR1c2VyX2lkXG4gICAgICB1cGRhdGVyIC0-PiBnYW1lX3N0YXRlOiBVcGRhdGUgZ2FtZSBzdGF0ZVxuICAgICAgdXBkYXRlciAtPj4gdXNlcl91cGRhdGVzOiBaQUREIHVzZXJfdXBkYXRlcyAkdXBkYXRlZF90aW1lc3RhbXAgJHVzZXJfaWRcbiAgICBlbmRcbiAgICBlbmRcbiIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In0sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoic2VxdWVuY2VEaWFncmFtXG4gICAgcGFydGljaXBhbnQgdXBkYXRlclxuICAgIHBhcnRpY2lwYW50IHVzZXJfdXBkYXRlc1xuICAgIHBhcnRpY2lwYW50IGdhbWVfc3RhdGVcblxuICAgIGxvb3BcbiAgICB1cGRhdGVyIC0-PiB1c2VyX3VwZGF0ZXM6IFpSQU5HRSB1c2VyX3VwZGF0ZXMgMCAwIFdJVEhTQ09SRVNcbiAgICB1c2VyX3VwZGF0ZXMgLT4-IHVwZGF0ZXI6ICgkdGltZXN0YW1wLCAkdXNlcl9pZClcbiAgICBhbHQgdGltZXN0YW1wIDwgbm93XG4gICAgICB1cGRhdGVyIC0-PiB1c2VyX3VwZGF0ZXM6IFpSRU0gdXNlcl91cGRhdGVzICR1c2VyX2lkXG4gICAgICB1cGRhdGVyIC0-PiBnYW1lX3N0YXRlOiBVcGRhdGUgZ2FtZSBzdGF0ZVxuICAgICAgdXBkYXRlciAtPj4gdXNlcl91cGRhdGVzOiBaQUREIHVzZXJfdXBkYXRlcyAkdXBkYXRlZF90aW1lc3RhbXAgJHVzZXJfaWRcbiAgICBlbmRcbiAgICBlbmRcbiIsIm1lcm1haWQiOnsidGhlbWUiOiJkZWZhdWx0In0sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)

### Client-Server Protocol

[Protocol Buffers](https://developers.google.com/protocol-buffers/) are used to define the messages sent between client/server and server/client. They are also as database models, i.e. how objects/documents are stored in Redis.

#### Client Messages

The following is not the exact definitions, they are here to describe on a higher-level what messages that exist and roughly the data they contain.

| Id  |  Type                 |  Payload           |
|-----|-----------------------|--------------------|
| ... |  "TAP"                | amount?            |
| ... |  "CONSTRUCT_BUILDING" | lotId, building    |
| ... |  "UPGRADE_BUILDING"   | lotId              |
| ... |  "TRAIN"              | education, amount  |
| ... |  "EXPAND"             |                    |
| ... |  "STEAL"              | amount, x, y       |

#### Server Messages

The following is not the exact definitions, they are here to describe on a higher-level what messages that exist and roughly the data they contain.

|  Type                 |  Payload           |
|-----------------------|--------------------|
|  "STATE_CHANGE"       | ...game_state      |
|  "RESPONSE"           | request_id, result |

## Tech Stack

- The Web App is/will be built with [TypeScript](https://www.typescriptlang.org/), _React_, [Zustand](https://github.com/pmndrs/zustand)
- The Web API and workers are/will be built with [Go](https://golang.org/) and [**Redis**](https://redis.io)
- Models for persistence and client/server messages are defined using [Protocol Buffers](https://developers.google.com/protocol-buffers/)

