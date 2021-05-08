#
# Builder
#
FROM golang:1.16-alpine3.13 AS builder

WORKDIR /app

# Install build tool dependencies
RUN apk add protobuf protobuf-dev make
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install modules
COPY go.mod go.sum /app/
RUN go mod download all

# Copy source and build all
COPY Makefile .
COPY protos/ protos/
COPY internal/ internal/
COPY cmd/ cmd/

RUN make build

#
# Web App Builder
#
FROM node:15 AS webapp-builder
WORKDIR /usr/src/app
COPY webapp/package.json webapp/package-lock.json ./
RUN npm ci
COPY webapp/ ./
COPY protos/ /usr/src/protos/
RUN npm run build


#
# Base runner
#
FROM alpine:3.13 AS base-runner

# Following commands are for installing CA certs (for proper functioning of HTTPS and other TLS)
RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

RUN adduser -D pizzatribes
USER pizzatribes
WORKDIR /home/pizzatribes
EXPOSE 8080

#
# API
#
FROM base-runner AS api
COPY --from=builder /app/out/pizza-tribes-api /app/pizza-tribes-api
CMD ["/app/pizza-tribes-api"]

#
# Updater
#
FROM base-runner AS updater
COPY --from=builder /app/out/pizza-tribes-updater /app/pizza-tribes-updater
CMD ["/app/pizza-tribes-updater"]

#
# Worker
#
FROM base-runner AS worker
COPY --from=builder /app/out/pizza-tribes-worker /app/pizza-tribes-worker
CMD ["/app/pizza-tribes-worker"]

#
# Migrator
#
FROM base-runner AS migrator
COPY --from=builder /app/out/pizza-tribes-migrator /app/pizza-tribes-migrator
CMD ["/app/pizza-tribes-migrator"]

#
# Web App
#
FROM nginx AS webapp
COPY --from=webapp-builder /usr/src/app/dist/ /usr/share/nginx/html

#
# Front
#
FROM caddy:2 AS front
COPY Caddyfile /etc/caddy/Caddyfile
