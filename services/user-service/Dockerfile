# Build stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod init streaming-service || true
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o service

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/service .
EXPOSE 8080
CMD ["./service"]
