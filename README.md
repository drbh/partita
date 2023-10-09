# partita

partita means "match" in Italian and this is a simple implementation of a match-making game server and a small snake like multiplayer game.

### 🎮 game

<a href="https://partita.fly.dev/" target="_blank">Try it out!</a>

### 🤳 Showcase

https://github.com/drbh/partita/assets/9896130/bde9db8c-4af2-4dab-8bb9-e9a40a92064f

### 🤓 nerdy notes

This implementation is based on the SOLID principles and more specifically an exercise of Dependency Injection using `wire` with `fiber` to build a easy to maintain and add new/extend features.

Since we depend on `wire` to build the application, we can easily split our application into multiple encapsulated modules and build them separately. Due to this all of the core services are split into directories at the top level of the project.

| Directory    | Description                                                       |
| ------------ | ----------------------------------------------------------------- |
| `background` | Background workers (update game state, match making, etc.)        |
| `collision`  | Collision detection (real number line based collision detection)  |
| `connection` | Connection management (dedicated cache for websocket connections) |
| `game`       | Game logic (includes game state and objects)                      |
| `match`      | Match making (simple match making based on player's elo)          |
| `redis`      | Redis client (mostly for match making)                            |
| `websocket`  | Websocket connection handling                                     |

### ☣️ disclaimer

This is a proof of concept and not a production ready implementation. There is much room for improvement and the code is not optimized for performance. The main goal of this project was to make a simple extendable game server that can be used to build a multiplayer games.
