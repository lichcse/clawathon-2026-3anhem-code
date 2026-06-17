# Voice Chat Room

**Team**: 3 anh em | **Group**: 3

A real-time voice chat application with seat-based room layout built with Go + React, PostgreSQL + Redis, and WebSocket.

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Start entire stack
docker compose up -d

# Open http://localhost:5173 in your browser
```

### Manual Setup

**Backend:**
```bash
cd backend && cp .env.example .env && go mod tidy && go run ./cmd/server/main.go
```

**Frontend:**
```bash
cd frontend && npm install && npm run dev
```

## Features

### Phase 1 — MVP
✅ User Registration & JWT Authentication
✅ Create/List Public Rooms
✅ 3×6 Seat Grid per Room
✅ Real-time Seat Occupancy (WebSocket)
✅ Microphone Mute/Unmute
✅ Backend-Mediated Audio Routing
✅ User Presence Tracking
✅ Cross-Server Redis Pub-Sub (horizontal scaling)

### Phase 2 — Scalability & Production Readiness
✅ Structured JSON logging
✅ Prometheus metrics (`/metrics`) — HTTP latency, WebSocket connections, audio frames
✅ Rate limiting — 100 req/min per user/IP (Redis sliding window)
✅ Panic recovery middleware
✅ `/health` liveness + `/ready` readiness probes (checks Postgres & Redis)
✅ Redis connection pooling with retry/backoff
✅ Graceful shutdown (SIGTERM/SIGINT with 10s drain)
✅ React `ErrorBoundary` with fallback UI
✅ Skeleton loading screens (seat grid, member list, room cards)
✅ Real-time audio level meter (mic visualization)
✅ User presence indicators (online dot per member)
✅ WebSocket reconnection toasts (reconnecting / reconnected / failed)
✅ Mobile-responsive layout (3-col grid on mobile, fixed controls bar)
✅ Kubernetes manifests (namespace, deployments, ingress, HPA, daily DB backup CronJob)

## Architecture

- **Backend**: Go (Gin, gorilla/websocket)
- **Frontend**: React 18 (Vite, TypeScript, Zustand, TailwindCSS)
- **Database**: PostgreSQL 15 (ACID, seat occupancy tracking)
- **Cache**: Redis 7 (pub-sub for events, session caching, rate limiting)
- **Audio**: Backend forwards binary WebSocket frames
- **Observability**: Prometheus metrics, structured JSON logs
- **Latency**: <100ms target

## API Endpoints

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
GET    /api/v1/auth/me

GET    /api/v1/rooms
POST   /api/v1/rooms
GET    /api/v1/rooms/:room_id
DELETE /api/v1/rooms/:room_id

POST   /api/v1/rooms/:room_id/seats/:seat_id/occupy
DELETE /api/v1/rooms/:room_id/seats/:seat_id

GET    /health      # liveness probe
GET    /ready       # readiness probe (checks Postgres + Redis)
GET    /metrics     # Prometheus metrics
```

## WebSocket Events

**Client → Server**: `join_room`, `sit_down`, `stand_up`, `mic_toggle`, `audio`, `leave_room`
**Server → Client**: `connected`, `room_snapshot`, `user_joined`, `seat_occupied`, `seat_vacated`, `user_mic_changed`, `user_left`, `audio`

## Project Structure

```
backend/              # Go server
├── cmd/server/       # Entry point (graceful shutdown)
├── internal/
│   ├── api/          # HTTP handlers & middleware (auth, rate limit, recovery)
│   ├── ws/           # WebSocket hub & handlers
│   ├── service/      # Business logic
│   ├── repository/   # Data access
│   ├── domain/       # Models & errors
│   └── infra/        # DB, Redis, config, Prometheus metrics

frontend/             # React app
├── src/
│   ├── pages/        # LoginPage, LandingPage, RoomPage
│   ├── components/   # UI components (ErrorBoundary, Skeleton, AudioMeter, PresenceIndicator)
│   ├── stores/       # Zustand stores
│   ├── services/     # API, WebSocket (reconnect callbacks), Audio
│   ├── types/        # TypeScript definitions
│   └── utils/        # cn helper

k8s/                  # Kubernetes manifests
├── namespace.yaml
├── configmap.yaml
├── secret.yaml
├── postgres-statefulset.yaml
├── redis-deployment.yaml
├── backend-deployment.yaml   # 2 replicas, rolling update, health probes
├── frontend-deployment.yaml
├── ingress.yaml              # WebSocket-aware
├── hpa.yaml                  # Autoscale backend 2→10
└── postgres-backup-cronjob.yaml  # Daily pg_dump at 2am UTC

docker-compose.yml    # Full local stack
```

## Testing

1. Register 2 users
2. User A creates a room
3. User B joins the room
4. Both sit on seats → see real-time updates
5. Toggle microphone → broadcast mute status + see audio level meter
6. Kill backend container → observe reconnection toasts
7. Leave room → cleanup

## Deploying to Kubernetes

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml      # update values first
kubectl apply -f k8s/
```
