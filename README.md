# ChatOrbit

## Build
```bash
docker build -t backend .
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