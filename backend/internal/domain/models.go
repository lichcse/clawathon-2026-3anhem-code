package domain

import (
	"errors"
	"time"
)

type User struct {
	ID               string    `json:"id"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	PasswordHash      string    `json:"-"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Room struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	MaxUsers    int       `json:"max_users"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Seat struct {
	ID            string     `json:"id"`
	RoomID        string     `json:"room_id"`
	RowNum        int        `json:"row"`
	ColNum        int        `json:"col"`
	OccupiedByID  *string    `json:"occupied_by_id"`
	OccupiedAt    *time.Time `json:"occupied_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type RoomMember struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	JoinedAt  time.Time `json:"joined_at"`
	LeftAt    *time.Time `json:"left_at"`
	IsMuted   bool      `json:"is_muted"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"-"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrRoomNotFound      = errors.New("room not found")
	ErrSeatNotFound      = errors.New("seat not found")
	ErrSeatOccupied      = errors.New("seat is already occupied")
	ErrRoomFull          = errors.New("room is full")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUsernameTaken     = errors.New("username is already taken")
	ErrEmailTaken        = errors.New("email is already taken")
)
