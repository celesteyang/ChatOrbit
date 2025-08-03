# ChatOrbit
A scalable real-time chat backend for live streams, built with WebSocket and Redis. Handles message broadcasting, rate limiting, and user bans.

## Core Features
- WebSocket Chat: Real-time messaging via Redis pub/sub
- Ban System: Admins can ban users with Redis keys + TTL
- Message Logging: All chat messages saved in MongoDB
- Presence Tracking: Redis Sets + TTL to track users in rooms
- Rate Limiting: Prevent spam using Redis Sorted Sets or Lua

## Tech Stack
- Go & Gin
- Redis (Pub/Sub, Lua, Sets)
- WebSocket (Gorilla)
- MongoDB
- JWT
- Docker & Docker Compose

## Architecture
Microservices managed with docker-compose.
- User: Profile & user queries
- Auth: User signup, login, JWT auth
- Chat: WebSocket server with Redis pub/sub
- Gateway: Unified API entrypoint

## Project Structure

```
.
├── .devcontainer/
│   ├── devcontainer.json
│   └── docker-compose.yaml
├── pkg/
│   ├── config/
│   ├── logger/
│   └── models/
├── services/
│   ├── auth/
│   ├── chat/
│   ├── gateway/
│   └── user/
├── docker-compose.yaml
├── Dockerfile
├── README.md
```

## Build
```bash
docker compose -f docker-compose.services.yaml up --build
```

## Related Project
[ChatOrbit Frontend](https://github.com/celesteyang/orbit-nexus-chat)
React app with auto-generated API clients via Orval.


## Running

### The entire application
To run the entire application, make sure that all services are not commented out in the `docker-compose.services.yaml` and run:
```bash
docker compoes -f docker-compose.services.yaml up --build
```
The `--build` is important to rebuild any services that might have changed since the last time we ran this command.

To develop a service independently from the other services, it can be beneficial for fast development to be able to directly run the service with `go`. To do this, comment out the relevant service in the `docker-compose.services.yaml` file.

### The Auth service
For example, for the `auth` service, we will run directly with `go` inside the dev-container as follows:
```bash
MONGO_URL=mongodb://localhost:27019 PORT=8089 go run .
```

## MongoDB

### CLI
To access the MongoDB records from the CLI, follow these steps:
1. First find the container corresponding to the MongoDB service and identify the container ID:
```bash
docker ps
```
2. Run `mongosh` inside the container with:
```bash
docker exec -it <container-ID> mongosh
```
This should run successfully.
3. Now run:
```bash
> show dbs
```
and identify you db name, e.g. `chatorbit`.
4. Now change the active DB:
```bash
use chatorbit
```
5. Now list the available collections with:
```bash
show collections
```
And identify the collection you want to inspect.
6. Finally, run:
```bash
db.<my_collection>.find().pretty()
```

### GUI
For convenience, there is a `Mongo-Express` container running alongside MongoDB that exposes a web-based UI for MongoDB.
To access it, simply run:
```bash
http://localhost:8091/
```
where `8091` is the port at which the `Mongo-Express` container is running. This is defined in `docker-compose.services.yaml`.