FROM golang:1.24 AS base
RUN apt-get update && apt-get install -y \
    vim \
    curl \
    tree \
    iputils-ping \
    sudo \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.5
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest

# Install mongosh in Debian-based golang image
RUN apt-get update && apt-get install -y curl gnupg \
 && curl -fsSL https://pgp.mongodb.com/server-7.0.asc | gpg --dearmor -o /usr/share/keyrings/mongodb-server-7.0.gpg \
 && echo "deb [ signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/debian bullseye/mongodb-org/7.0 main" > /etc/apt/sources.list.d/mongodb-org-7.0.list \
 && apt-get update && apt-get install -y mongodb-mongosh \
 && rm -rf /var/lib/apt/lists/*

RUN useradd -m -s /bin/bash -G sudo vscode \
    && echo 'vscode ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

ENV GOPATH=/home/vscode/go
ENV PATH=$PATH:/home/vscode/go/bin
RUN mkdir -p /home/vscode/go && chown -R vscode:vscode /home/vscode/go

WORKDIR /workspace