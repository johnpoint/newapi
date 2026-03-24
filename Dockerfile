# syntax=docker/dockerfile:1

# ---- builder ----
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o newapi .

# ---- final ----
FROM alpine:latest

RUN apk add --no-cache ca-certificates

# /workspace is the mount point for the user's Go project
WORKDIR /workspace

COPY --from=builder /build/newapi /usr/local/bin/newapi

ENTRYPOINT ["newapi"]
