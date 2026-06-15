# Voice Chat Room - Phase 1 MVP

**Team**: 3 anh em | **Group**: 3

A real-time voice chat application with seat-based room layout built with Go + React, PostgreSQL + Redis, and WebSocket.

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Start entire stack
docker-compose up -d

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

## MVP Features

✅ User Registration & JWT Authentication  
✅ Create/List Public Rooms  
✅ 3x6 Seat Grid per Room  
✅ Real-time Seat Occupancy (WebSocket)  
✅ Microphone Mute/Unmute  
✅ Backend-Mediated Audio Routing  
✅ User Presence Tracking  
✅ Cross-Server Redis Pub-Sub (for horizontal scaling)

## Architecture

- **Backend**: Go (Gin framework, gorilla/websocket)
- **Frontend**: React 18 (Vite, TypeScript, Zustand, TailwindCSS)
- **Database**: PostgreSQL (ACID, seat occupancy tracking)
- **Cache**: Redis (pub-sub for events, session caching)
- **Audio**: Backend forwards binary WebSocket frames (no codec needed in MVP)
- **Latency**: <100ms target (real-time optimization)

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
```

## WebSocket Events

**Client → Server**: join_room, sit_down, stand_up, mic_toggle, audio, leave_room  
**Server → Client**: connected, room_snapshot, user_joined, seat_occupied, seat_vacated, user_mic_changed, user_left, audio

## Project Structure

```
backend/              # Go server
├── cmd/server/       # Entry point
├── internal/
│   ├── api/          # HTTP handlers & middleware
│   ├── ws/           # WebSocket hub & handlers
│   ├── service/      # Business logic
│   ├── repository/   # Data access
│   ├── domain/       # Models & errors
│   └── infra/        # DB, Redis, config

frontend/             # React app
├── src/
│   ├── pages/        # LoginPage, LandingPage, RoomPage
│   ├── components/   # UI components
│   ├── stores/       # Zustand stores
│   ├── services/     # API, WebSocket, Audio
│   ├── types/        # TypeScript definitions
│   └── styles/       # TailwindCSS

docker-compose.yml    # Full stack setup
```

## Testing

1. Register 2 users
2. User A creates a room
3. User B joins the room
4. Both sit on seats → see real-time updates
5. Toggle microphone → broadcast mute status
6. Leave room → cleanup

## Next Phase

- WebRTC P2P (reduce bandwidth)
- Chat messaging
- Speaker queue system
- Audience mode
- Audio recording
- Room persistence
- Screen sharing