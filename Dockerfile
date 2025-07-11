# Stage 1: Build
FROM golang:1.22 AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    npm \
    libgtk-3-dev \
    libwebkit2gtk-4.0-dev \
    build-essential \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install Wails CLI and build app
RUN go install github.com/wailsapp/wails/v2/cmd/wails@latest
RUN wails build -production

# Stage 2: Runtime
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    libgtk-3-0 \
    libwebkit2gtk-4.0-37 \
    libjavascriptcoregtk-4.0-18 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/build/bin/* /app/vatsim-tracker
COPY --from=builder /app/branding.json /app/branding.json
COPY --from=builder /app/branding /app/branding

ENTRYPOINT ["/app/vatsim-tracker"]
