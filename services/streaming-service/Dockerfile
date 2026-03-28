# Build stage
FROM golang:1.24.5 AS builder 
WORKDIR /app

# 1. Copy only module files first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# 2. Copy the rest of the source code
COPY . .

# 3. Explicitly build (Change "." to your subfolder path if main.go isn't in root)
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/api

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/service .
EXPOSE 8080
CMD ["./service"]