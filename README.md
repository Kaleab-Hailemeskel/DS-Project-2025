# Meta Vibes Backend

Distributed music platform backend composed of Go microservices, Nginx gateway, PostgreSQL, Redis, Kafka, and Kubernetes manifests. This file covers only the backend (Compose and K8s).

## What’s Inside

- Services: `services/` contains Go services (user, song, streaming). Compose also wires in playlist-service and adaptive-engine (images built from sibling paths if present).
- Gateway: [api-gateway/nginx.conf](api-gateway/nginx.conf) reverse proxy for songs and streaming; user route commented for now.
- Local orchestration: [docker/docker-compose-orginal.yml](docker/docker-compose-orginal.yml) full stack; [docker/docker-compose.yml](docker/docker-compose.yml) slim stack.
- Data init: [docker/postgres/init.sql](docker/postgres/init.sql) seeds users, playlists, songs tables on first Postgres start.
- Kubernetes: [k8s/](k8s) contains ingress, Postgres clusters (CloudNativePG), Redis StatefulSets, and per-service manifests.
- Media: `SONG_ARCHIVE/` holds uploaded song assets.

## Architecture Snapshot

- API gateway: Nginx routes `/api/songs/*` to song-service and `/api/stream/*` to streaming-service; health at `/health`.
- Services: Go HTTP APIs exposing users, songs (CRUD + upload), and streaming. Redis used for caching, Kafka for events, Postgres for persistence.
- Infra: Postgres 16, Redis 7, Kafka + Zookeeper, pgAdmin, Kafdrop included in the full Compose stack.

## Run Locally with Docker Compose (full stack)

From the docker folder run the full environment (infra + all services + UIs):

```bash
cd Backend/DS-Project-2025/docker
docker compose -f docker-compose-orginal.yml up --build
```

Service endpoints

- Gateway: http://localhost:8000
- User API: http://localhost:8001
- Playlist API: http://localhost:8002
- Song API: http://localhost:8003
- Adaptive Engine: http://localhost:8004
- Streaming API: http://localhost:8005
- pgAdmin: http://localhost:5050 (admin@admin.com / admin)
- Kafdrop: http://localhost:9000
- Postgres: localhost:5432 (user postgres, password postgres)

Stop and clean:

```bash
docker compose -f docker-compose-orginal.yml down -v
```

Notes

- Postgres schema auto-seeds from [docker/postgres/init.sql](docker/postgres/init.sql) on first boot; data persists in the `postgres_data` volume.
- Health checks guard Postgres, Redis, Kafka; wait until they are healthy before hitting APIs.
- The slim [docker/docker-compose.yml](docker/docker-compose.yml) starts only services without infra; use it when you already have external databases/caches.

### API Gateway routing

- `/api/songs/*` → song-service (upload route increases `client_max_body_size` to 20 MB).
- `/api/stream/*` → streaming-service.
- `/health` → gateway 200 OK.
- User route is present but commented; enable when user-service is deployed.

### Database tables

Created at startup: `users`, `playlists`, `songs` (UUID primary key) as defined in [docker/postgres/init.sql](docker/postgres/init.sql).

## Kubernetes Manifests (k8s/)

- Ingress: [k8s/nginx/ingress.yaml](k8s/nginx/ingress.yaml) routes `song-app.test` host to song-api, streaming-api, and user-api services with CORS and 25 MB body limit.
- Postgres: [k8s/postgres-sql/cluster.yaml](k8s/postgres-sql/cluster.yaml) provisions two CloudNativePG clusters (`song-postgres` and `user-postgres`, 2 instances each, 5Gi storage) bootstrapped from secrets.
- Redis: [k8s/redis-cache/song-redis-statefulset.yaml](k8s/redis-cache/song-redis-statefulset.yaml) defines a 2-replica StatefulSet with master election and a single service endpoint.
- Services: service-specific folders (song-service, streaming-service, user-service) contain Deployment/Service/ConfigMap manifests; update image tags and secrets before applying.
- Secrets: raw and sealed secrets templates exist under `k8s/postgres-sql` and `k8s/redis-cache`; replace placeholder values for production.

Apply example (after configuring images/secrets and pointing DNS for `song-app.test`):

```bash
kubectl apply -f k8s/postgres-sql/
kubectl apply -f k8s/redis-cache/
kubectl apply -f k8s/song-service/
kubectl apply -f k8s/streaming-service/
kubectl apply -f k8s/user-service/
kubectl apply -f k8s/nginx/ingress.yaml
```

## Troubleshooting

- Ports in use: stop previous containers (`docker ps`) or change mappings in [docker/docker-compose-orginal.yml](docker/docker-compose-orginal.yml).
- Services not healthy: `docker compose -f docker-compose-orginal.yml logs -f <service>`.
- Fresh DB seed: remove the `postgres_data` volume, then rerun the stack.
- Ingress unreachable: confirm hosts file/DNS maps `song-app.test` to cluster ingress IP and that the ingress class matches your controller.
