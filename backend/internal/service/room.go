package service

import (
	"fmt"

	"github.com/google/uuid"

	"voicechat/internal/domain"
	"voicechat/internal/repository"
)

type RoomService struct {
	roomRepo   *repository.RoomRepository
	seatRepo   *repository.SeatRepository
	memberRepo *repository.RoomMemberRepository
}

func NewRoomService(roomRepo *repository.RoomRepository, seatRepo *repository.SeatRepository, memberRepo *repository.RoomMemberRepository) *RoomService {
	return &RoomService{
		roomRepo:   roomRepo,
		seatRepo:   seatRepo,
		memberRepo: memberRepo,
	}
}

type CreateRoomRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoomDTO struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	OwnerID     string      `json:"owner_id"`
	MaxUsers    int         `json:"max_users"`
	IsPublic    bool        `json:"is_public"`
	MemberCount int         `json:"member_count"`
	Seats       []SeatDTO   `json:"seats"`
	Members     []MemberDTO `json:"members"`
}

type SeatDTO struct {
	ID           string `json:"id"`
	Row          int    `json:"row"`
	Col          int    `json:"col"`
	OccupiedByID *string `json:"occupied_by_id"`
}

type MemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsMuted  bool   `json:"is_muted"`
}

func (s *RoomService) CreateRoom(ownerID string, req *CreateRoomRequest) (*RoomDTO, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("room name is required")
	}

	room := &domain.Room{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		MaxUsers:    18,
		IsPublic:    true,
	}

	err := s.roomRepo.Create(room)
	if err != nil {
		return nil, err
	}

	// Create 3x6 seat grid
	err = s.seatRepo.CreateSeatsForRoom(room.ID, 3, 6)
	if err != nil {
		return nil, err
	}

	return s.GetRoom(room.ID)
}

func (s *RoomService) GetRoom(roomID string) (*RoomDTO, error) {
	room, err := s.roomRepo.GetByID(roomID)
	if err != nil {
		return nil, err
	}

	seats, err := s.seatRepo.GetSeatsForRoom(roomID)
	if err != nil {
		return nil, err
	}

	members, err := s.memberRepo.GetMembersWithUsername(roomID)
	if err != nil {
		return nil, err
	}

	// Build username lookup for seat enrichment
	usernameMap := make(map[string]string, len(members))
	for _, m := range members {
		usernameMap[m.UserID] = m.Username
	}

	seatDTOs := make([]SeatDTO, len(seats))
	for i, seat := range seats {
		seatDTOs[i] = SeatDTO{
			ID:           seat.ID,
			Row:          seat.RowNum,
			Col:          seat.ColNum,
			OccupiedByID: seat.OccupiedByID,
		}
	}

	memberDTOs := make([]MemberDTO, len(members))
	for i, member := range members {
		memberDTOs[i] = MemberDTO{
			UserID:   member.UserID,
			Username: member.Username,
			IsMuted:  member.IsMuted,
		}
	}

	return &RoomDTO{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		OwnerID:     room.OwnerID,
		MaxUsers:    room.MaxUsers,
		IsPublic:    room.IsPublic,
		MemberCount: len(members),
		Seats:       seatDTOs,
		Members:     memberDTOs,
	}, nil
}

func (s *RoomService) ListRooms(limit, offset int) ([]RoomDTO, error) {
	rooms, err := s.roomRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	var dtos []RoomDTO
	for _, room := range rooms {
		count, _ := s.roomRepo.GetMemberCount(room.ID)
		dto := RoomDTO{
			ID:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			OwnerID:     room.OwnerID,
			MaxUsers:    room.MaxUsers,
			IsPublic:    room.IsPublic,
			MemberCount: count,
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

func (s *RoomService) DeleteRoom(roomID, userID string) error {
	room, err := s.roomRepo.GetByID(roomID)
	if err != nil {
		return err
	}

	if room.OwnerID != userID {
		return domain.ErrUnauthorized
	}

	return s.roomRepo.Delete(roomID)
}

func (s *RoomService) JoinRoom(roomID, userID string) error {
	room, err := s.roomRepo.GetByID(roomID)
	if err != nil {
		return err
	}

	count, err := s.roomRepo.GetMemberCount(roomID)
	if err != nil {
		return err
	}

	if count >= room.MaxUsers {
		return domain.ErrRoomFull
	}

	return s.memberRepo.Join(roomID, userID)
}

func (s *RoomService) LeaveRoom(roomID, userID string) error {
	// Vacate any seats user is sitting on
	_ = s.seatRepo.VacateSeatsByUser(roomID, userID)
	return s.memberRepo.Leave(roomID, userID)
}
