# ================================================================================
FROM golang:1.23 AS builder-backend

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o main

# ================================================================================
FROM node:20 AS builder-frontend

WORKDIR /frontend

COPY ./frontend ./frontend

WORKDIR /frontend/frontend

RUN npm install

RUN npm run build

# ================================================================================
FROM alpine:3.18 AS builder-maa-cli

ARG maa_cli_ver="v0.5.4" # 默认值
ENV MAA_CLI_VER=${maa_cli_ver}

RUN apk add --no-cache \
    wget \
    tar \
    bash

WORKDIR /builder-maa-cli

RUN arch=$(uname -m) && \
    if [ "$arch" = "x86_64" ]; then \
        file="maa_cli-x86_64-unknown-linux-gnu.tar.gz"; \
    elif [ "$arch" = "aarch64" ]; then \
        file="maa_cli-aarch64-unknown-linux-gnu.tar.gz"; \
    else \
        echo "Unsupported architecture: $arch" && exit 1; \
    fi && \
    wget https://github.com/MaaAssistantArknights/maa-cli/releases/download/${MAA_CLI_VER}/$file && \
    tar -xzf $file --strip-components=1 && \
    rm $file

# ================================================================================
FROM debian:latest

RUN apt-get update && apt-get install -y ca-certificates git adb && rm -rf /var/lib/apt/lists/*

ENV LANG=C.UTF-8

ENV MAA_CONFIG_DIR=/app/config/maa

ENV PROXY=

ENV CLIENT_TYPE=Bilibili

WORKDIR /app

COPY --from=builder-backend /app/main /app/main
COPY --from=builder-frontend /frontend/frontend/dist /app/dist
COPY --from=builder-maa-cli /builder-maa-cli/maa /usr/bin/maa

COPY ./config/templateConfig.json /tmp/config/config.json

COPY ./config/maa /tmp/config/maa

VOLUME /app/config

COPY ./entrypoint.sh /app/entrypoint.sh

CMD ["./entrypoint.sh"]
