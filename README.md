# ChatOrbit
A scalable real-time chat backend for live streams, built with WebSocket and Redis. Handles message broadcasting, rate limiting, and user bans.

## Architecture
Microservices using Go + Gin + MongoDB:

auth: User signup, login, JWT auth

user: Profile & user queries

chat: WebSocket server with Redis pub/sub

gateway: Unified API entrypoint

Managed with docker-compose

## Core Features
WebSocket Chat: Real-time messaging via Redis pub/sub

Ban System: Admins can ban users with Redis keys + TTL

Message Logging: All chat messages saved in MongoDB

Presence Tracking: Redis Sets + TTL to track users in rooms

Rate Limiting: Prevent spam using Redis Sorted Sets or Lua

## Build
```bash
docker compose -f docker-compose.services.yaml up --build
```

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