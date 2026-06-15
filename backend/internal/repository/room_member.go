package repository

import (
	"database/sql"
	"time"

	"voicechat/internal/domain"
)

type RoomMemberRepository struct {
	db *sql.DB
}

func NewRoomMemberRepository(db *sql.DB) *RoomMemberRepository {
	return &RoomMemberRepository{db: db}
}

func (r *RoomMemberRepository) Join(roomID, userID string) error {
	query := `
		INSERT INTO room_members (id, room_id, user_id, joined_at)
		VALUES (gen_random_uuid(), $1, $2, $3)
		ON CONFLICT (room_id, user_id) DO UPDATE
		SET left_at = NULL
	`

	_, err := r.db.Exec(query, roomID, userID, time.Now())
	return err
}

func (r *RoomMemberRepository) Leave(roomID, userID string) error {
	query := `
		UPDATE room_members
		SET left_at = $1
		WHERE room_id = $2 AND user_id = $3
	`

	_, err := r.db.Exec(query, time.Now(), roomID, userID)
	return err
}

func (r *RoomMemberRepository) GetMembers(roomID string) ([]domain.RoomMember, error) {
	query := `
		SELECT id, room_id, user_id, joined_at, left_at, is_muted
		FROM room_members
		WHERE room_id = $1 AND left_at IS NULL
	`

	rows, err := r.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.RoomMember
	for rows.Next() {
		var member domain.RoomMember
		err := rows.Scan(&member.ID, &member.RoomID, &member.UserID, &member.JoinedAt, &member.LeftAt, &member.IsMuted)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, rows.Err()
}

func (r *RoomMemberRepository) SetMuteStatus(roomID, userID string, isMuted bool) error {
	query := `
		UPDATE room_members
		SET is_muted = $1
		WHERE room_id = $2 AND user_id = $3
	`

	_, err := r.db.Exec(query, isMuted, roomID, userID)
	return err
}
