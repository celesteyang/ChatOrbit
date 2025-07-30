FROM golang:1.24 AS base
RUN apt-get update && apt-get install -y \
    vim \
    curl \
    tree \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.5

RUN useradd -m vscode
WORKDIR /workspace