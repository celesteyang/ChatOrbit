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
To run the entire application, make sure that all services are not commented out in the `docker-compose.services.yaml` and run (make sure to run this in a terminal outside of the dev container because docker is no installed in the dev container):
```bash
docker compose -f docker-compose.services.yaml up --build
```
The `--build` is important to rebuild any services that might have changed since the last time we ran this command.

To develop a service independently from the other services, it can be beneficial for fast development to be able to directly run the service with `go`. To do this, comment out the relevant service in the `docker-compose.services.yaml` file.

### The Auth service

For example, for the `auth` service, we will run directly with `go` inside the dev-container as follows:
```bash
JWT_SECRET="your_secret_key" MONGO_URL=mongodb://host.docker.internal:27019 PORT=8089 go run .
```
> Note: When running  **inside the dev container**, `localhost` refers to the dev container itself, **not the host machine**.  
If MongoDB or Redis is running in a Docker container (via Docker Compose) and mapped to a host port (e.g., `27019` for MongoDB), we must use `host.docker.internal` to connect from the dev container instead of `localhost`.

#### Example API Usage
**Swagger** 
http://localhost:8089/swagger/index.html

**Register Test**
> Note: These are manual test commands, not using Swagger.

```bash
curl -X POST http://localhost:8089/register \
  -H "Content-Type: application/json" \
  -d '{"email":"abc@example.com", "username":"test", "password":"12345678"}'
```

**Login Test**
```bash
curl -X POST http://localhost:8089/login \
  -H "Content-Type: application/json" \
  -d '{"email":"abc@example.com", "password":"12345678"}'
```

### The Chat Service

For example, for the `chat` service, we will run directly with `go` inside the dev-container as follows:
```bash
JWT_SECRET="your_secret_key" MONGO_URL="mongodb://host.docker.internal:27019" REDIS_ADDR="host.docker.internal:6381" PORT=8088 go run .
```
#### Creating a Room
Rooms are created on demand via a simple REST call. This is useful when the frontend navigates to a room like `music` before
anyone has joined it.

```bash
curl -X POST http://localhost:8088/chat/rooms \
  -H "Content-Type: application/json" \
  -d '{"room_id": "music"}'
```

The same call is idempotent—if the room already exists it simply returns the requested `room_id`.

#### WebSocket Testing with `wscat`
After logging in with the auth service and getting a JWT, you can test the WebSocket connection with `wscat`. Pass `room_id` in
the WebSocket URL to join a room (defaults to `general`). If the room does not exist yet, it will be created automatically.
```bash
npm install -g wscat
```
```bash
wscat -c "ws://localhost:8088/ws/chat?token=<YOUR_JWT_TOKEN>&room_id=my-room"
```
**Verifying Real-Time Broadcast with Multiple Terminals**
To confirm that real-time messaging is working correctly, we need to simulate multiple users. This tests the WebSocket connections and the Redis Pub/Sub broadcast functionality.

1. Obtain a Second JWT Token
Use auth service to get a unique JWT for a second user.
2. Establish Connections
Open two separate terminal windows and connect to the chat service, each with a different user's JWT.
```bash
# Terminal 1: User 1's connection
wscat -c "ws://localhost:8088/ws/chat?token=<USER_1_JWT_TOKEN>&room_id=my-room"
```
```bash
# Terminal 2: User 2's connection
wscat -c "ws://localhost:8088/ws/chat?token=<USER_2_JWT_TOKEN>&room_id=my-room"
```
3. Test Real-Time Communication

```bash
# Terminal 1: User 1
# Send a message:
{"room_id": "my-room", "content": "Hello!"}
```
```bash
# Terminal 2: User 2
# Send a reply:
{"room_id": "my-room", "content": "Hi there! Got your message instantly."}
```
Both terminals should display messages in real time:
```json
{"id":"000000000000000000000000","room_id":"general","user_id":"690766ae876dd929bee54fcd","content":"Hello!","timestamp":"2025-11-04T13:21:51.729935846Z"}
"id":"000000000000000000000000","room_id":"general","user_id":"69075d4b876dd929bee54fcc","content":"Hi there! Got your message instantly.","timestamp":"2025-11-04T13:22:11.273811914Z"}
```
## Forwarded Ports in Dev Containers

When we run the services inside a **VS Code dev container**, the ports the services listen on (like `8088` for chat) are **inside the container**, not directly on the host machine.  
VS Code automatically forwards these ports to the host, but:

- **The forwarded port on the host may not always match the internal port.**
- For example, the chat service might listen on `8088` inside the container, but VS Code could forward it to `localhost:37347` or another random port on the host.

### How to Find the Correct Forwarded Port

1. **Open the "Ports" panel in VS Code** (bottom panel or via Command Palette: "Ports: Focus on Ports View").
2. **Look for the port that the service is listening on** (e.g., `8088`).  
   The "Local Address" column shows the forwarded address on the host (e.g., `localhost:37347`).

### Example

If the "Ports" panel shows:

| Port | Forwarded Address |
|------|-------------------|
| 8088 | localhost:37347   |

We should connect with:

```bash
wscat -c "ws://localhost:37347/ws/chat?token=<YOUR_JWT_TOKEN>"
```

**Not**:

```bash
wscat -c "ws://localhost:8088/ws/chat?token=<YOUR_JWT_TOKEN>"
```
(unless the local address is actually `localhost:8088`)

### How to Always Use the Same Port

To make sure VS Code always forwards, for example, `8088` to `localhost:8088`,  
add this to `.devcontainer/devcontainer.json`:

```json
{
  "forwardPorts": [8088, 8089]
}
```
Then **rebuild or reload** the dev container.

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