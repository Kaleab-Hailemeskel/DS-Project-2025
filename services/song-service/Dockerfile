# Build stage
FROM golang:1.24.5 AS builder
WORKDIR /app

# 1. Copy go.mod and go.sum first. 
# This allows Docker to cache your dependencies so they don't 
# re-download every time you change a line of code.
COPY go.mod go.sum* ./
RUN go mod download

# 2. Copy the rest of the source code
COPY . .

# 3. Build the application
# We use "." to build the package in the current directory.
# If your main.go is in a folder like cmd/api, change "." to "./cmd/api"
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/api

# Run stage
FROM alpine:latest
# Add ca-certificates in case your service calls external APIs (HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/service .

# Optional: Ensure the binary is executable
RUN chmod +x /app/service
# Install ffmpeg for audio processing
RUN apk add --no-cache ffmpeg

EXPOSE 8080

CMD ["./service"]