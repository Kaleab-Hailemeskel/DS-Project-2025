## Purpose
Help AI coding agents become productive quickly in this repository by documenting the big-picture architecture, developer workflows, service patterns, and concrete examples taken from the codebase.

## **Architecture Overview**
- **Microservices (Go)**: Five small services live under `services/` (`user-service`, `playlist-service`, `song-service`, `adaptive-engine`, `streaming-service`). Each service currently exposes a minimal HTTP API with a `/health` endpoint implemented in `services/<service>/main.go`.
- **API Gateway**: `api-gateway/` contains an Nginx config used as the single entry point (reverse proxy). It maps external ports to the internal service endpoints.
- **Infrastructure (Docker Compose)**: `docker/docker-compose.yml` wires up Postgres, Redis, Kafka+Zookeeper, pgAdmin, Kafdrop, and the services. Use the `docker` folder to start the full stack.

## **Run / Developer Workflows**
- Full stack (first run):
  - `cd docker`
  - `docker compose up --build` (first time; ~2 minutes)
- Full stack (detached):
  - `cd docker`
  - `docker compose up -d`
- Stop and remove volumes:
  - `cd docker`
  - `docker compose down -v`
- Health quick checks (example):
  - `curl http://localhost:8000/health` (API gateway)
  - `curl http://localhost:8000/api/users/health` (proxied to user service)

## **Service Patterns & Conventions**
- All services listen on port `8080` inside their container (`http.ListenAndServe(":8080", ...)` in `main.go`). The API gateway maps each service to a different public port.
- Dockerfiles perform a multi-stage Go build:
  - `go mod init streaming-service || true`
  - `go mod tidy`
  - `CGO_ENABLED=0 GOOS=linux go build -o service`
  - Result is copied into a slim `alpine` image and run.
  Note: `go mod init` inside Docker means no explicit `go.mod` file lives in the repo — be careful when running services locally.
- Environment variables (set in Compose):
  - `POSTGRES_URL` (postgres connection)
  - `REDIS_URL` (redis connection)
  - `KAFKA_BROKERS` (kafka brokers)
  - `SERVICE_PORT` (services default to internal 8080)

## **Important Files to Reference**
- `README.md` — project overview and run commands.
- `docker/docker-compose.yml` — full stack wiring and environment variables.
- `api-gateway/nginx.conf` — routing and port mappings (source of truth for external endpoints).
- `services/*/main.go` and `services/*/Dockerfile` — canonical service behavior and container build steps.
- `docker/postgres/init.sql` — DB schema auto-creation used by Compose.

## **Concrete Examples / Quick Edits**
- To add a new HTTP endpoint to a service: edit `services/<service>/main.go` and add `http.HandleFunc("/your-route", handler)`. Rebuild the image via `cd docker && docker compose build <service>` and `docker compose up -d`.
- To run a single service locally for quick iteration (recommended):
  - Open a PowerShell terminal in `services/<service>` and run `go run .` (the binary listens on `:8080`). If `go` complains about module missing, run `go mod init <module>` locally and `go mod tidy` — this mirrors the Dockerfile behavior.

## **Integration Points & External Tools**
- Postgres: `localhost:5432` (user: `postgres`, pass: `postgres`). Schema lives in `docker/postgres/init.sql`.
- pgAdmin UI: `http://localhost:5050` (login `admin@admin.com` / `admin`).
- Kafdrop (Kafka UI): `http://localhost:9000`.

## **What to Avoid / Notes for an AI Agent**
- Do not assume business logic exists beyond health endpoints — many services are scaffolds. Check `main.go` before changing APIs.
- The Dockerfiles initialize Go modules inside the build stage; avoid relying on a repository-level `go.mod` unless you add and commit it explicitly.

## **When to Ask the Human**
- If you need the intended public API paths (beyond `/health`) — ask which service owns the route and expected request/response.
- If you plan to change database schemas — confirm migration strategy and whether Compose-managed Postgres should be reused or replaced.

If any section is unclear or you'd like additional examples (e.g., example PR diff for a new endpoint, or a local debug recipe with `delve`), tell me which area to expand.
