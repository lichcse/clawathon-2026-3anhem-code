package repository

import (
	"database/sql"

	"voicechat/internal/domain"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *domain.Room) error {
	query := `
		INSERT INTO rooms (id, name, description, owner_id, max_users, is_public)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, room.ID, room.Name, room.Description, room.OwnerID, room.MaxUsers, room.IsPublic).
		Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt)

	return err
}

func (r *RoomRepository) GetByID(id string) (*domain.Room, error) {
	room := &domain.Room{}
	query := `
		SELECT id, name, description, owner_id, max_users, is_public, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&room.ID, &room.Name, &room.Description, &room.OwnerID,
		&room.MaxUsers, &room.IsPublic, &room.CreatedAt, &room.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrRoomNotFound
	}
	return room, err
}

func (r *RoomRepository) List(limit, offset int) ([]domain.Room, error) {
	query := `
		SELECT id, name, description, owner_id, max_users, is_public, created_at, updated_at
		FROM rooms
		WHERE is_public = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []domain.Room
	for rows.Next() {
		var room domain.Room
		err := rows.Scan(
			&room.ID, &room.Name, &room.Description, &room.OwnerID,
			&room.MaxUsers, &room.IsPublic, &room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, rows.Err()
}

func (r *RoomRepository) Delete(id string) error {
	query := `DELETE FROM rooms WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrRoomNotFound
	}

	return nil
}

func (r *RoomRepository) GetMemberCount(roomID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM room_members WHERE room_id = $1 AND left_at IS NULL`
	err := r.db.QueryRow(query, roomID).Scan(&count)
	return count, err
}
