# DS-Project-2025
### A music Streaming Service
This is the GitHub README for a full-featured music streaming service designed to deliver a high-quality, uninterrupted audio experience. The platform is built for performance and scale, focusing on robust content delivery and a wide selection of music.

### ✨ Core Features
- Adaptive Audio Streaming: Enjoy crystal-clear audio quality that seamlessly adjusts to your network conditions for uninterrupted listening. This feature ensures playback remains smooth even as bandwidth fluctuates.

- Global Content Delivery (CDN): We leverage a global Content Delivery Network (CDN) to distribute music files closer to every user. This results in fast load times and reliable playback, regardless of your geographical location.

- Search and Discovery: Powerful search functionality allows users to quickly find the artists, songs, or albums they are looking for.

- User Profiles & Social Features: Users can manage their listening data, follow other accounts, and share music with friends.


<img width="472" height="453" alt="image" src="https://github.com/user-attachments/assets/6f578580-628a-40b7-9fe6-0a129f0daf2e" />

text## Features Implemented
- Nginx reverse proxy (single entry point)
- PostgreSQL with auto-initialized tables
- Redis caching
- Kafka + Zookeeper event bus
- pgAdmin (DB GUI) – http://localhost:5050
- Kafdrop (Kafka GUI) – http://localhost:9000
- Health checks on every service
- Docker Compose one-command deployment

## How to Run (Team Instructions)

### 1. Start everything (first time – takes ~2 minutes)
```bash
cd docker
docker compose up --build
2. Future restarts
Bashcd docker
docker compose up -d
3. Stop everything
Bashdocker compose down -v
Service URLs (via API Gateway

ServiceEndpointDirect PortAPI Gatewayhttp://localhost:80008000User Servicehttp://localhost:8000/api/users8001Playlist Servicehttp://localhost:8000/api/playlists8002Song Servicehttp://localhost:8000/api/songs8003Adaptive Enginehttp://localhost:8000/api/adaptive8004Streaming Servicehttp://localhost:8000/api/stream8005
Quick health checks
Bashcurl http://localhost:8000/health
curl http://localhost:8000/api/users/health
curl http://localhost:8000/api/playlists/health
# ... etc
Developer Tools

pgAdmin → http://localhost:5050
login: admin@admin.com / admin
Kafdrop (Kafka UI) → http://localhost:9000
PostgreSQL direct: localhost:5432 (user: postgres, pass: postgres)

Database Tables (auto-created)
SQLusers      (id, username, email, password_hash, created_at)
playlists  (id, user_id → users.id, name, created_at)
songs      (id, title, artist, duration, file_path, uploaded_at)
Environment Variables (already set in compose)
All services automatically receive:
textPOSTGRES_URL=postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable
REDIS_URL=redis://redis:6379
KAFKA_BROKERS=kafka:9092
SERVICE_PORT=8080
