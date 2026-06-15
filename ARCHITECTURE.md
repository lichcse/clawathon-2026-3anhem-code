# Voice Chat Room - System Architecture

## 1. Overall Architecture

### High-Level Design
```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (React)                         │
│  - Room UI (seats, users, controls)                              │
│  - Audio capture/playback (WebAudio API)                         │
│  - WebSocket real-time updates                                   │
└──────────────────────┬──────────────────────────────────────────┘
                       │ WebSocket
                       │ HTTP REST API
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    API Gateway / Load Balancer                   │
│              (NGINX or cloud load balancer)                      │
└──────────────────────┬──────────────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│  Go Server 1 │ │  Go Server 2 │ │  Go Server N │
│  (WebSocket) │ │  (WebSocket) │ │  (WebSocket) │
│  (HTTP API)  │ │  (HTTP API)  │ │  (HTTP API)  │
└──────────────┘ └──────────────┘ └──────────────┘
        │              │              │
        └──────────────┼──────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
        ▼              ▼              ▼
    ┌────────┐   ┌─────────┐   ┌──────────┐
    │ Redis  │   │Postgres │   │S3/Blob   │
    │Pub/Sub │   │State DB │   │Audio Rec.│
    │Caching │   │Users    │   │          │
    │Sessions│   │Rooms    │   │          │
    │        │   │Seats    │   │          │
    └────────┘   └─────────┘   └──────────┘
```

### Key Design Principles

1. **Stateless Server Layer**: Go servers are horizontally scalable; shared state lives in Redis/Postgres
2. **Event-Driven Communication**: WebSocket events for real-time updates, pub-sub for cross-server coordination
3. **Backend-Mediated Audio**: Simplifies recording, moderation, and state sync (trades off bandwidth for control)
4. **Bounded Contexts**: Room, User, Seat, Audio, Notification services (can be separate microservices later)
5. **< 100ms Latency**: Optimized with:
   - Direct WebSocket (no HTTP polling)
   - Binary protocol for audio frames (optional protobuf)
   - Redis for fast state queries
   - Minimal CPU on audio forwarding (avoid codec operations in MVP)

---

## 2. Proposed Technologies

### Backend Stack
| Component | Choice | Rationale |
|-----------|--------|-----------|
| **Language** | Go 1.22+ | Goroutines excel at concurrent connections; < 1ms scheduling latency |
| **Framework** | Gin (HTTP) + gorilla/websocket (WS) | Lightweight, high throughput, minimal latency |
| **Database** | PostgreSQL 15+ | ACID guarantees for room/user state; JSONB for flexible config |
| **Cache/Pub-Sub** | Redis 7+ | Fast pub-sub for cross-server events; session caching |
| **Audio Processing** | Native Go + WebAudio API (browser) | Opus codec in browser, raw frames via WS (codec cost deferred) |
| **Deployment** | Kubernetes (or simple systemd on AWS/GCP) | Auto-scaling, health checks, rolling updates |

### Frontend Stack
| Component | Choice | Rationale |
|-----------|--------|-----------|
| **Framework** | React 18+ | Component-driven, good for real-time updates with hooks |
| **State** | Zustand (or Jotai) | Lightweight, minimal boilerplate vs Redux |
| **WebSocket** | Socket.IO + native fallback | Auto-reconnect, built-in namespacing, dev tools |
| **Audio Capture** | WebAudio API + getUserMedia | Browser standard, cross-platform |
| **UI** | TailwindCSS + Radix UI | Accessible, fast iteration |
| **Testing** | Vitest + React Testing Library | Fast, good DX |

### Infrastructure
| Component | Choice | Rationale |
|-----------|--------|-----------|
| **Container** | Docker | Reproducible deployments, standard for Go |
| **Orchestration** | Kubernetes or Cloud Run | Horizontal scaling, managed alternative |
| **Storage** | PostgreSQL backups to S3; S3 for recordings | Durable, cost-effective |
| **Monitoring** | Prometheus (metrics) + Grafana (dashboards) | Open-source, scalable observability |
| **Logging** | Structured logging (JSON) to ELK/CloudLogging | Centralized, searchable logs |

---

## 3. Tradeoffs

### Backend-Mediated Audio vs Peer-to-Peer

| Aspect | Backend-Mediated | Peer-to-Peer (WebRTC) |
|--------|------------------|----------------------|
| **Bandwidth** | High (N streams to server) | Low (direct user-to-user) |
| **Latency** | Slightly higher | Slightly lower |
| **Recording** | Trivial (server has all streams) | Complex (SFU/MCU needed) |
| **Moderation** | Easy (audio passes through) | Impossible without SFU |
| **Development** | Simpler | Complex signaling, NAT traversal |
| **MVP fit** | ✓ Choose for MVP | ✗ Defer to v2 |

**Decision**: Start with backend-mediated. Migrate to WebRTC SFU in v2 if bandwidth becomes a bottleneck.

### Synchronous REST vs Event-Driven Pub-Sub

| Aspect | REST (Polling/Sync) | Pub-Sub Events |
|--------|---------------------|-----------------|
| **Latency** | 100-500ms (polling) | < 50ms (pushed) |
| **Simplicity** | Higher | Medium |
| **Cross-server sync** | Requires eventual consistency | Built-in |

**Decision**: Use WebSocket events for real-time updates; REST only for non-time-sensitive operations (auth, metadata).

### Single Database vs Multiple Services

| Aspect | Single Postgres | Microservices |
|--------|-----------------|---------------|
| **Consistency** | ACID transactions | Eventual consistency |
| **Scaling** | Vertical until ~10K QPS | Easier horizontal |
| **Complexity** | Low | High |
| **Deployment** | Simple | Complex |

**Decision**: Single Postgres + Redis pub-sub for MVP (scales to 10K QPS). Split into microservices (User, Room, Audio) in v2 if needed.

---

## 4. Bounded Contexts / Modules

### Backend Bounded Contexts

```
┌────────────────────────────────────────────────────────────────┐
│ Room Service                                                    │
│ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │
│ │ Room Repo    │ │ Room Logic   │ │ Room Events  │             │
│ │ - CRUD       │ │ - Seat logic │ │ - Created    │             │
│ │ - Validation │ │ - Max users  │ │ - Closed     │             │
│ └──────────────┘ └──────────────┘ └──────────────┘             │
└────────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────────┐
│ User Service                                                    │
│ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │
│ │ User Repo    │ │ User Logic   │ │ User Events  │             │
│ │ - CRUD       │ │ - Profile    │ │ - Joined     │             │
│ │ - Sessions   │ │ - Auth state │ │ - Left       │             │
│ └──────────────┘ └──────────────┘ └──────────────┘             │
└────────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────────┐
│ Seat Service                                                    │
│ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │
│ │ Seat Repo    │ │ Seat Logic   │ │ Seat Events  │             │
│ │ - CRUD       │ │ - Occupancy  │ │ - Occupied   │             │
│ │ - State      │ │ - Validation │ │ - Vacated    │             │
│ └──────────────┘ └──────────────┘ └──────────────┘             │
└────────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────────┐
│ Audio Service                                                   │
│ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │
│ │ Audio Router │ │ Stream Mgmt  │ │ Audio Events │             │
│ │ - Forward    │ │ - Track users│ │ - Started    │             │
│ │ - Mix (later)│ │ - Mute state │ │ - Stopped    │             │
│ └──────────────┘ └──────────────┘ └──────────────┘             │
└────────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────────┐
│ WebSocket Handler (Coordinator)                                │
│ - Delegates to services                                        │
│ - Manages client connections                                   │
│ - Broadcasts events                                            │
└────────────────────────────────────────────────────────────────┘
```

### Frontend Bounded Contexts

```
┌─────────────────────────────────────────────────────────────┐
│ Pages                                                        │
│ ┌──────────────────────┐ ┌──────────────────────┐            │
│ │ RoomPage             │ │ LandingPage          │            │
│ │ - Main room UI       │ │ - Room creation      │            │
│ └──────────────────────┘ └──────────────────────┘            │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ Components (Presentational)                                  │
│ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐           │
│ │ SeatGrid     │ │ UserList     │ │ Controls     │           │
│ │ AudioDisplay │ │ Chat         │ │ Settings     │           │
│ └──────────────┘ └──────────────┘ └──────────────┘           │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ Stores (State Management via Zustand)                        │
│ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐           │
│ │ useRoomStore │ │ useUserStore │ │ useAudioStore│           │
│ │ - Room state │ │ - Local user │ │ - Mic/volume│           │
│ │ - Seats      │ │ - Profiles   │ │ - Streams   │           │
│ └──────────────┘ └──────────────┘ └──────────────┘           │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ Services (Business Logic)                                    │
│ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐           │
│ │ WebSocket    │ │ Audio        │ │ API          │           │
│ │ - Connection │ │ - Capture    │ │ - REST calls │           │
│ │ - Events     │ │ - Playback   │ │ - Auth       │           │
│ └──────────────┘ └──────────────┘ └──────────────┘           │
└─────────────────────────────────────────────────────────────┘
```

---

## 5. Backend APIs

### REST Endpoints (Non-Real-Time Operations)

```
Authentication
  POST   /api/v1/auth/register              # Create user account
  POST   /api/v1/auth/login                 # Issue session token
  POST   /api/v1/auth/logout                # Invalidate session
  GET    /api/v1/auth/me                    # Get current user

Rooms (CRUD)
  GET    /api/v1/rooms                      # List rooms (paginated)
  POST   /api/v1/rooms                      # Create room
  GET    /api/v1/rooms/:room_id             # Get room details
  PUT    /api/v1/rooms/:room_id             # Update room settings
  DELETE /api/v1/rooms/:room_id             # Delete room (owner only)

Room Members
  GET    /api/v1/rooms/:room_id/members     # List members in room
  DELETE /api/v1/rooms/:room_id/members/:user_id  # Remove user

Seats
  GET    /api/v1/rooms/:room_id/seats       # Get seat layout + occupancy
  POST   /api/v1/rooms/:room_id/seats/:seat_id/occupy  # Claim seat
  DELETE /api/v1/rooms/:room_id/seats/:seat_id         # Leave seat

User Profile
  GET    /api/v1/users/:user_id             # Get user profile
  PUT    /api/v1/users/me                   # Update own profile
```

### Response Format

```json
{
  "status": "success|error",
  "data": { /* payload */ },
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

---

## 6. WebSocket Events

### Connection Lifecycle

```
Client connects
  → server: authentication token in query param or header
  → server: validates token, creates connection record in Redis
  → client: receives "connected" event with current state snapshot

Client joins room
  → client: emits "join_room" { room_id, user_id }
  → server: validates user, adds to room pub-sub channel
  → server: emits "user_joined" to all in room (broadcast)
  → client: receives full room state snapshot

Client sits on seat
  → client: emits "sit_down" { room_id, seat_id }
  → server: acquires lock, validates seat empty, updates DB
  → server: emits "seat_occupied" to room (broadcast)
  → server: emits "user_position_changed" (optional, for analytics)

Client starts audio
  → client: starts capturing audio → sends frames via binary WS
  → server: routes frames to all other users in room
  → server: forwards audio frames to users (binary messages)

Client disconnects
  → server: cleans up WebSocket, publishes "user_left" event
  → all clients: receive "user_left" { user_id }
```

### Event Schema

#### Client → Server

```js
{
  "event": "join_room",
  "payload": {
    "room_id": "uuid",
    "user_id": "uuid"
  }
}

{
  "event": "sit_down",
  "payload": {
    "room_id": "uuid",
    "seat_id": "seat_1_1",
    "user_id": "uuid"
  }
}

{
  "event": "stand_up",
  "payload": {
    "room_id": "uuid",
    "seat_id": "seat_1_1",
    "user_id": "uuid"
  }
}

{
  "event": "mic_toggle",
  "payload": {
    "room_id": "uuid",
    "user_id": "uuid",
    "is_muted": true|false
  }
}

{
  "event": "leave_room",
  "payload": {
    "room_id": "uuid",
    "user_id": "uuid"
  }
}

// Binary audio frame
[0x01, ...audio_bytes]  // First byte = frame type, rest = audio payload
```

#### Server → Client

```js
{
  "event": "connected",
  "payload": {
    "user_id": "uuid",
    "server_time": 1718550000000
  }
}

{
  "event": "room_snapshot",
  "payload": {
    "room_id": "uuid",
    "name": "Main Hall",
    "seats": [
      { "id": "seat_1_1", "row": 1, "col": 1, "occupied_by": "uuid" },
      { "id": "seat_1_2", "row": 1, "col": 2, "occupied_by": null },
      ...
    ],
    "members": [
      { "user_id": "uuid", "username": "alice", "seat_id": "seat_1_1", "is_muted": false },
      ...
    ]
  }
}

{
  "event": "user_joined",
  "payload": {
    "user_id": "uuid",
    "username": "bob",
    "joined_at": 1718550001000
  }
}

{
  "event": "seat_occupied",
  "payload": {
    "room_id": "uuid",
    "seat_id": "seat_1_1",
    "user_id": "uuid",
    "username": "alice"
  }
}

{
  "event": "seat_vacated",
  "payload": {
    "room_id": "uuid",
    "seat_id": "seat_1_1",
    "user_id": "uuid"
  }
}

{
  "event": "user_mic_changed",
  "payload": {
    "user_id": "uuid",
    "is_muted": true|false
  }
}

{
  "event": "user_left",
  "payload": {
    "user_id": "uuid",
    "left_at": 1718550005000
  }
}

// Binary audio frame
[0x01, frame_number_4bytes, user_id_16bytes, audio_bytes...]
```

### Cross-Server Coordination (Redis Pub-Sub)

```
Channels:
  - room:{room_id}:events         # Room-specific events
  - room:{room_id}:audio          # Audio frames for the room
  - global:user_presence          # User online/offline events (optional)
  - global:notifications          # System notifications
```

---

## 7. Database Schema

### PostgreSQL Tables

```sql
-- Users
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  profile_picture_url TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Rooms
CREATE TABLE rooms (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  description TEXT,
  owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
  max_users INT DEFAULT 18,  -- 3 rows * 6 columns
  is_public BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  
  INDEX idx_owner_id (owner_id),
  INDEX idx_is_public (is_public)
);

-- Seats (fixed grid: 3 rows, 6 cols per room)
CREATE TABLE seats (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  row_num INT NOT NULL,        -- 1, 2, 3
  col_num INT NOT NULL,        -- 1, 2, 3, 4, 5, 6
  occupied_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
  occupied_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  
  UNIQUE (room_id, row_num, col_num),
  INDEX idx_room_occupied (room_id, occupied_by_id)
);

-- Room Membership (tracks who has joined a room, presence)
CREATE TABLE room_members (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  joined_at TIMESTAMP DEFAULT NOW(),
  left_at TIMESTAMP,
  is_muted BOOLEAN DEFAULT FALSE,
  
  UNIQUE (room_id, user_id),
  INDEX idx_room_user (room_id, user_id),
  INDEX idx_user_room (user_id, room_id)
);

-- Audio Recordings (future feature)
CREATE TABLE audio_recordings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  started_by_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
  s3_uri VARCHAR(512) NOT NULL,
  duration_seconds INT,
  started_at TIMESTAMP DEFAULT NOW(),
  ended_at TIMESTAMP,
  
  INDEX idx_room_id (room_id),
  INDEX idx_started_by (started_by_id)
);

-- Sessions (for auth token tracking)
CREATE TABLE sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash VARCHAR(255) NOT NULL UNIQUE,
  ip_address INET,
  user_agent TEXT,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  
  INDEX idx_user_sessions (user_id),
  INDEX idx_expires_at (expires_at)
);

-- Audit Log (for debugging, analytics)
CREATE TABLE audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_type VARCHAR(100) NOT NULL,  -- user_joined, user_left, seat_occupied, etc
  room_id UUID REFERENCES rooms(id) ON DELETE SET NULL,
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  data JSONB,  -- event payload
  created_at TIMESTAMP DEFAULT NOW(),
  
  INDEX idx_room_event (room_id, event_type),
  INDEX idx_user_event (user_id, event_type),
  INDEX idx_created_at (created_at)
);
```

### Redis Schema

```
Session/Connection State (short-lived, TTL 30min):
  ws:session:{user_id} = { 
    conn_id, 
    current_room_id, 
    joined_at, 
    server_instance_id 
  }

User Presence (TTL 30sec, updated on heartbeat):
  presence:{room_id}:{user_id} = { 
    username, 
    seat_id, 
    is_muted, 
    last_heartbeat 
  }

Room State Cache (TTL 5min, invalidated on changes):
  room:{room_id}:cache = {
    name,
    description,
    seats: [{ id, row, col, occupied_by }],
    members: [{ user_id, username, seat_id, is_muted }]
  }

Audio Stream Metadata (ephemeral, TTL 5min):
  audio:stream:{room_id}:{user_id} = {
    bitrate,
    codec,
    frame_rate,
    last_frame_at
  }
```

---

## 8. Frontend Architecture

### Component Tree

```
App
├── Router
│   ├── LandingPage
│   │   └── CreateRoomModal
│   │   └── RoomList
│   ├── RoomPage
│   │   ├── SeatGrid
│   │   │   └── Seat (interactive)
│   │   │   │   └── Avatar
│   │   │   │   └── Username
│   │   ├── RoomHeader
│   │   │   └── RoomTitle
│   │   │   └── MemberCount
│   │   │   └── LeaveButton
│   │   ├── UserPanel
│   │   │   └── UserList
│   │   │   │   └── UserCard (+ mute button)
│   │   │   └── MyProfileCard
│   │   ├── ControlPanel
│   │   │   └── MicrophoneToggle
│   │   │   └── VolumeSlider
│   │   │   └── SettingsButton
│   │   └── NotificationArea
│   │       └── Toast notifications
│   └── NotFoundPage
├── AuthProvider
├── WebSocketProvider
└── ErrorBoundary
```

### State Management (Zustand)

```ts
// Room Store
useRoomStore: {
  roomId: string | null
  roomName: string
  seats: Seat[]
  members: User[]
  
  joinRoom(roomId)
  leaveRoom()
  sitDown(seatId)
  standUp()
  setSeatOccupied(seatId, userId)
  setSeatVacated(seatId)
  updateMemberMuteStatus(userId, isMuted)
}

// Audio Store
useAudioStore: {
  isMicOn: boolean
  volume: number
  activeStreams: Map<userId, AudioStreamData>
  
  toggleMic()
  setVolume(level)
  addStream(userId, stream)
  removeStream(userId)
  getMixedAudio()
}

// Auth Store
useAuthStore: {
  user: User | null
  isAuthenticated: boolean
  token: string | null
  
  login(email, password)
  logout()
  register(email, username, password)
  setUser(user)
}

// UI Store
useUIStore: {
  isMuted: boolean
  showSettings: boolean
  notifications: Notification[]
  
  toggleMute()
  showNotification(message, type)
  clearNotification(id)
}
```

### WebSocket Service

```ts
class WebSocketService {
  private ws: WebSocket
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000
  
  connect(token: string, roomId?: string)
  disconnect()
  emit(event: string, payload: any)
  on(event: string, handler: (payload) => void)
  off(event: string, handler: (payload) => void)
  
  // Internal
  private attemptReconnect()
  private handleMessage(event: MessageEvent)
  private handleBinaryMessage(data: ArrayBuffer)  // For audio
}
```

### Audio Service

```ts
class AudioService {
  private localStream: MediaStream
  private audioContext: AudioContext
  private remoteStreams: Map<userId, AudioStreamData>
  
  // Capture
  async startCapture()
  async stopCapture()
  async setMicMuted(isMuted: boolean)
  
  // Playback
  addRemoteStream(userId: string, stream: Uint8Array)
  removeRemoteStream(userId: string)
  setRemoteVolume(userId: string, volume: number)
  
  // Utility
  getLocalVolume(): number
  getRemoteVolume(userId: string): number
  async testMicrophone(): Promise<boolean>
}
```

### API Service

```ts
class ApiService {
  private baseUrl: string
  private token: string
  
  // Auth
  login(email, password): Promise<{ token, user }>
  register(username, email, password): Promise<{ token, user }>
  
  // Rooms
  listRooms(page?, limit?): Promise<Room[]>
  createRoom(name, description): Promise<Room>
  getRoom(roomId): Promise<Room>
  deleteRoom(roomId): Promise<void>
  
  // Seats
  getSeatLayout(roomId): Promise<Seat[]>
  occupySeat(roomId, seatId): Promise<Seat>
  vacateSeat(roomId, seatId): Promise<Seat>
  
  // Members
  getRoomMembers(roomId): Promise<User[]>
  removeFromRoom(roomId, userId): Promise<void>
}
```

---

## 9. Folder Structure

```
.
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler/          # HTTP handlers
│   │   │   │   ├── room.go
│   │   │   │   ├── user.go
│   │   │   │   ├── seat.go
│   │   │   │   ├── auth.go
│   │   │   │   └── health.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── cors.go
│   │   │   │   └── logging.go
│   │   │   └── router.go
│   │   ├── ws/                   # WebSocket
│   │   │   ├── handler.go        # WS event handlers
│   │   │   ├── hub.go            # Connection manager
│   │   │   ├── message.go        # Message types
│   │   │   └── broadcaster.go    # Event broadcasting
│   │   ├── service/              # Business logic
│   │   │   ├── room.go
│   │   │   ├── user.go
│   │   │   ├── seat.go
│   │   │   ├── audio.go
│   │   │   └── auth.go
│   │   ├── repository/           # Data access
│   │   │   ├── room.go
│   │   │   ├── user.go
│   │   │   ├── seat.go
│   │   │   └── session.go
│   │   ├── domain/               # Domain models
│   │   │   ├── room.go
│   │   │   ├── user.go
│   │   │   ├── seat.go
│   │   │   └── error.go
│   │   ├── infra/
│   │   │   ├── db/               # Database setup
│   │   │   │   ├── postgres.go
│   │   │   │   └── migration/
│   │   │   │       ├── 001_create_users.up.sql
│   │   │   │       ├── 002_create_rooms.up.sql
│   │   │   │       └── ...
│   │   │   ├── cache/            # Redis setup
│   │   │   │   ├── redis.go
│   │   │   │   └── pubsub.go
│   │   │   └── config/
│   │   │       └── config.go
│   │   └── util/
│   │       ├── logger.go
│   │       ├── jwt.go
│   │       └── errors.go
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── .env.example
│   └── README.md
│
├── frontend/
│   ├── src/
│   │   ├── app.tsx               # Root component
│   │   ├── main.tsx              # Entry point
│   │   ├── pages/
│   │   │   ├── LandingPage.tsx
│   │   │   ├── RoomPage.tsx
│   │   │   └── NotFoundPage.tsx
│   │   ├── components/
│   │   │   ├── layout/
│   │   │   │   ├── Header.tsx
│   │   │   │   └── Footer.tsx
│   │   │   ├── room/
│   │   │   │   ├── SeatGrid.tsx
│   │   │   │   ├── Seat.tsx
│   │   │   │   ├── UserPanel.tsx
│   │   │   │   └── UserCard.tsx
│   │   │   ├── controls/
│   │   │   │   ├── MicToggle.tsx
│   │   │   │   ├── VolumeControl.tsx
│   │   │   │   └── SettingsPanel.tsx
│   │   │   └── common/
│   │   │       ├── Modal.tsx
│   │   │       ├── Button.tsx
│   │   │       └── Toast.tsx
│   │   ├── stores/
│   │   │   ├── useRoomStore.ts
│   │   │   ├── useAudioStore.ts
│   │   │   ├── useAuthStore.ts
│   │   │   └── useUIStore.ts
│   │   ├── services/
│   │   │   ├── api.ts
│   │   │   ├── websocket.ts
│   │   │   ├── audio.ts
│   │   │   └── storage.ts
│   │   ├── hooks/
│   │   │   ├── useWebSocket.ts
│   │   │   ├── useAudio.ts
│   │   │   ├── useAuth.ts
│   │   │   └── useRoom.ts
│   │   ├── types/
│   │   │   ├── index.ts
│   │   │   ├── room.ts
│   │   │   ├── user.ts
│   │   │   ├── events.ts
│   │   │   └── audio.ts
│   │   ├── styles/
│   │   │   ├── globals.css
│   │   │   └── theme.ts
│   │   └── utils/
│   │       ├── logger.ts
│   │       ├── validators.ts
│   │       └── helpers.ts
│   ├── public/
│   │   ├── index.html
│   │   └── favicon.ico
│   ├── tests/
│   │   ├── unit/
│   │   ├── integration/
│   │   └── e2e/
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── tailwind.config.ts
│   ├── package.json
│   └── README.md
│
├── docker-compose.yml            # Local dev: Go, React, Postgres, Redis
├── k8s/                          # Kubernetes manifests
│   ├── backend-deployment.yaml
│   ├── backend-service.yaml
│   ├── postgres-statefulset.yaml
│   ├── redis-deployment.yaml
│   └── ingress.yaml
├── .github/
│   ├── workflows/
│   │   ├── test.yml
│   │   ├── build.yml
│   │   └── deploy.yml
├── ARCHITECTURE.md               # This file
├── CONTRIBUTING.md
└── README.md
```

---

## 10. Implementation Roadmap

### Phase 1: MVP (Weeks 1-3)
**Goal**: Basic voice room with 3x6 seat grid

- [ ] **Backend**
  - [x] Project setup (Go, gin, gorilla/websocket)
  - [ ] Database schema (users, rooms, seats, room_members)
  - [ ] Authentication (JWT, sessions)
  - [ ] REST API (rooms CRUD, seats CRUD, members)
  - [ ] WebSocket handler (connect, join, sit, stand, leave)
  - [ ] Audio routing (forward binary frames to room)
  - [ ] Redis setup (pub-sub for cross-server events)

- [ ] **Frontend**
  - [ ] React app setup (vite, React Router)
  - [ ] Landing page (room list, create room form)
  - [ ] Room page (seat grid, user list, controls)
  - [ ] Zustand stores (room, audio, auth, ui)
  - [ ] WebSocket service (connect, emit, handle events)
  - [ ] Audio capture/playback (getUserMedia, WebAudio API)
  - [ ] UI/UX (TailwindCSS, responsive design)

- [ ] **Testing & Deployment**
  - [ ] Unit tests (Go services, React components)
  - [ ] Integration tests (API + DB, WebSocket)
  - [ ] Docker setup (backend, frontend, postgres, redis)
  - [ ] Manual testing (latency, audio quality, seat occupation)

### Phase 2: Scalability & Polish (Weeks 4-5)
**Goal**: Horizontal scaling, error handling, monitoring

- [ ] **Backend**
  - [ ] Connection pooling (postgres, redis)
  - [ ] Load testing (1K concurrent users)
  - [ ] Reconnect logic (graceful recovery)
  - [ ] Logging & observability (structured logs, Prometheus metrics)
  - [ ] Rate limiting (prevent abuse)
  - [ ] Error handling (validation, timeouts, edge cases)

- [ ] **Frontend**
  - [ ] Error boundaries and fallback UI
  - [ ] Loading states (skeleton screens)
  - [ ] Reconnection UI (toast notifications)
  - [ ] Audio level meters (input/output)
  - [ ] User presence indicators
  - [ ] Mobile responsiveness

- [ ] **Deployment**
  - [ ] Kubernetes manifests (backend, postgres, redis)
  - [ ] Health checks & readiness probes
  - [ ] Rolling updates strategy
  - [ ] Database backups

### Phase 3: Core Features (Weeks 6-8)
**Goal**: Multiple rooms, basic moderation, improvements

- [ ] **Backend**
  - [ ] Room metadata (description, tags, capacity limits)
  - [ ] User profiles (avatar, bio, preferences)
  - [ ] Room permissions (owner, members, public/private)
  - [ ] Mute/unmute events and validation
  - [ ] Connection timeout handling (30s inactivity → auto-leave)

- [ ] **Frontend**
  - [ ] Profile editing (avatar upload, bio)
  - [ ] Room settings modal (rename, capacity, visibility)
  - [ ] Member management (remove, mute, kick)
  - [ ] User preferences (auto-gain control, noise suppression)
  - [ ] Dark mode support

### Phase 4: Advanced Features (Weeks 9-12)
**Goal**: Recording, analytics, enhanced experience

- [ ] **Audio Recording**
  - [ ] Server-side recording (write audio frames to file)
  - [ ] S3 storage integration
  - [ ] Playback service (serve recordings)
  - [ ] Timestamp markers (who spoke when)

- [ ] **Notifications**
  - [ ] User joined/left alerts
  - [ ] Missed messages (if chat added)
  - [ ] Email digest (daily summary)

- [ ] **Analytics & Debugging**
  - [ ] Grafana dashboard (concurrent users, latency, packet loss)
  - [ ] User analytics (session duration, room activity)
  - [ ] Error tracking (Sentry)

### Phase 5: Roadmap for Future Versions

#### v2 (Post-MVP)
- [ ] **WebRTC Peer-to-Peer** (reduce bandwidth)
  - Replace backend audio routing with SFU (Selective Forwarding Unit)
  - Go library: pion/webrtc
  - Frontend: webrtc-adapter, local offer/answer generation

- [ ] **Chat & Reactions**
  - Message storage in postgres
  - Emoji reactions (stored in room_events table)
  - Typing indicators via WebSocket

- [ ] **Speaker Queue**
  - Users request to speak
  - Owner/moderator approves
  - Queue visualized in UI

- [ ] **Audience Mode**
  - Read-only for non-speakers
  - Broadcast to 1000s of listeners (separate SFU/streaming setup)

- [ ] **Room Persistence**
  - Recurring rooms (daily standups)
  - Auto-archive recordings by room

#### v3+
- [ ] **Screen Sharing** (WebRTC + optional server relay)
- [ ] **Live Streaming** (RTMP to YouTube, Twitch)
- [ ] **Private Messaging** (DM between users)
- [ ] **Advanced Moderation** (shadowban, message filtering)
- [ ] **Mobile App** (React Native or Flutter)
- [ ] **Admin Dashboard** (user/room management, analytics)
- [ ] **API Rate Limiting & Quota** (tiers: free, pro, enterprise)
- [ ] **Webhooks** (integrations: Slack, Discord)

---

## Performance & Latency Targets

| Metric | Target | How |
|--------|--------|-----|
| WebSocket connect | < 50ms | Direct connection, minimal TLS overhead |
| Room join response | < 100ms | Cache room state in Redis, fast snapshot send |
| Seat occupation | < 100ms | Direct DB update, publish event immediately |
| Audio frame latency | < 100ms | Frame size 20ms, network 50-80ms |
| Presence update | < 100ms | Direct publish via Redis pub-sub |
| Reconnect | < 2s | Attempt 3x with 500ms backoff |

## Security Considerations

1. **Authentication**: JWT + session tokens, refresh token rotation
2. **Authorization**: Room owner/moderator checks, row-level security
3. **Rate Limiting**: 100 reqs/min per user, 10 joins/min per room
4. **Validation**: Input sanitization, XSS protection (React), CORS
5. **Data**: TLS for all connections, password hashing (bcrypt)
6. **Audio**: No encryption required in MVP (can add DTLS later)

## Testing Strategy

1. **Unit Tests**: Services, repositories, handlers (>80% coverage)
2. **Integration Tests**: API + DB, WebSocket events, pub-sub
3. **Load Tests**: k6 or Apache JMeter, simulate 1K concurrent users
4. **E2E Tests**: Playwright, happy path + error scenarios
5. **Manual Testing**: Latency with network throttling, audio quality

---

## Next Steps

1. **Review & Align**: Share this doc, gather feedback on tech choices & approach
2. **Spike**: Prototype audio routing in Go + React (1 day)
3. **Setup**: Initialize projects, DB migrations, Docker Compose
4. **Implement Phase 1**: MVP features in priority order
5. **Iterate**: Weekly demos, gather user feedback, refine roadmap
