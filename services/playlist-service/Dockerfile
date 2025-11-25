# Build stage
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod init playlist-service
RUN go build -o service .

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/service .
EXPOSE 8080
CMD ["./service"]
