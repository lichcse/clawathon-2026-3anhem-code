package repository

import (
	"database/sql"
	"time"

	"voicechat/internal/domain"
)

type SeatRepository struct {
	db *sql.DB
}

func NewSeatRepository(db *sql.DB) *SeatRepository {
	return &SeatRepository{db: db}
}

func (r *SeatRepository) CreateSeatsForRoom(roomID string, rows, cols int) error {
	for row := 1; row <= rows; row++ {
		for col := 1; col <= cols; col++ {
			query := `
				INSERT INTO seats (id, room_id, row_num, col_num)
				VALUES (gen_random_uuid(), $1, $2, $3)
			`
			_, err := r.db.Exec(query, roomID, row, col)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *SeatRepository) GetSeatsForRoom(roomID string) ([]domain.Seat, error) {
	query := `
		SELECT id, room_id, row_num, col_num, occupied_by_id, occupied_at, created_at
		FROM seats
		WHERE room_id = $1
		ORDER BY row_num, col_num
	`

	rows, err := r.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []domain.Seat
	for rows.Next() {
		var seat domain.Seat
		err := rows.Scan(
			&seat.ID, &seat.RoomID, &seat.RowNum, &seat.ColNum,
			&seat.OccupiedByID, &seat.OccupiedAt, &seat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	return seats, rows.Err()
}

func (r *SeatRepository) GetSeatByID(seatID string) (*domain.Seat, error) {
	seat := &domain.Seat{}
	query := `
		SELECT id, room_id, row_num, col_num, occupied_by_id, occupied_at, created_at
		FROM seats
		WHERE id = $1
	`

	err := r.db.QueryRow(query, seatID).Scan(
		&seat.ID, &seat.RoomID, &seat.RowNum, &seat.ColNum,
		&seat.OccupiedByID, &seat.OccupiedAt, &seat.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrSeatNotFound
	}
	return seat, err
}

func (r *SeatRepository) OccupySeat(seatID, userID string) error {
	now := time.Now()
	// Allow same user to re-occupy (idempotent for double-click) but block others
	query := `
		UPDATE seats
		SET occupied_by_id = $1, occupied_at = $2
		WHERE id = $3 AND (occupied_by_id IS NULL OR occupied_by_id = $1)
	`

	result, err := r.db.Exec(query, userID, now, seatID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrSeatOccupied
	}

	return nil
}

func (r *SeatRepository) VacateSeat(seatID string) error {
	query := `
		UPDATE seats
		SET occupied_by_id = NULL, occupied_at = NULL
		WHERE id = $1
	`

	result, err := r.db.Exec(query, seatID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrSeatNotFound
	}

	return nil
}

func (r *SeatRepository) VacateSeatsByUser(roomID, userID string) error {
	query := `
		UPDATE seats
		SET occupied_by_id = NULL, occupied_at = NULL
		WHERE room_id = $1 AND occupied_by_id = $2
	`

	_, err := r.db.Exec(query, roomID, userID)
	return err
}
