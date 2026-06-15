package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createUsersTable,
		createRoomsTable,
		createSeatsTable,
		createRoomMembersTable,
		createSessionsTable,
		createAuditLogsTable,
	}

	for i, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("failed to run migration %d: %w", i, err)
		}
	}

	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	username VARCHAR(255) UNIQUE NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	profile_picture_url TEXT,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
`

const createRoomsTable = `
CREATE TABLE IF NOT EXISTS rooms (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
	max_users INT DEFAULT 18,
	is_public BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rooms_owner_id ON rooms(owner_id);
CREATE INDEX IF NOT EXISTS idx_rooms_is_public ON rooms(is_public);
`

const createSeatsTable = `
CREATE TABLE IF NOT EXISTS seats (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
	row_num INT NOT NULL,
	col_num INT NOT NULL,
	occupied_by_id UUID REFERENCES users(id) ON DELETE SET NULL,
	occupied_at TIMESTAMP,
	created_at TIMESTAMP DEFAULT NOW(),

	UNIQUE(room_id, row_num, col_num)
);

CREATE INDEX IF NOT EXISTS idx_seats_room_occupied ON seats(room_id, occupied_by_id);
`

const createRoomMembersTable = `
CREATE TABLE IF NOT EXISTS room_members (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	joined_at TIMESTAMP DEFAULT NOW(),
	left_at TIMESTAMP,
	is_muted BOOLEAN DEFAULT FALSE,

	UNIQUE(room_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_room_members_room_user ON room_members(room_id, user_id);
CREATE INDEX IF NOT EXISTS idx_room_members_user_room ON room_members(user_id, room_id);
`

const createSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	token_hash VARCHAR(255) NOT NULL UNIQUE,
	ip_address INET,
	user_agent TEXT,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
`

const createAuditLogsTable = `
CREATE TABLE IF NOT EXISTS audit_logs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	event_type VARCHAR(100) NOT NULL,
	room_id UUID REFERENCES rooms(id) ON DELETE SET NULL,
	user_id UUID REFERENCES users(id) ON DELETE SET NULL,
	data JSONB,
	created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_room_event ON audit_logs(room_id, event_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_event ON audit_logs(user_id, event_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
`
