# User Service (PoC)

Development PoC for user registration, login, and profile access.

Required env vars (for local/dev):

- `POSTGRES_URL` (e.g. `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`)
- `JWT_SECRET` (dev-only secret in `docker/docker-compose.yml`: `dev-secret-change-me`)
- `JWT_EXPIRY_MINUTES` (optional, default `60`)
- `BCRYPT_COST` (optional, default `12`)
- `SERVICE_PORT` (optional, default `8080`)

Run locally (fast iteration):

PowerShell example:

```powershell
$env:POSTGRES_URL = 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable'
$env:JWT_SECRET = 'dev-secret-change-me'
$env:JWT_EXPIRY_MINUTES = '60'
cd services\user-service
go run .
```

When running the full stack via Docker Compose the API gateway exposes the endpoints at `http://localhost:8000/api/users/*`.
