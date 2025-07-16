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
