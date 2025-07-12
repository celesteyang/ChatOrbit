FROM golang:1.24 AS base
RUN apt-get update && apt-get install -y \
    vim \
    curl \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -m vscode
WORKDIR /workspace