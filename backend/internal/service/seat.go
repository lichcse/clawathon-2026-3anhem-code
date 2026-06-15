package service

import (
	"voicechat/internal/domain"
	"voicechat/internal/repository"
)

type SeatService struct {
	seatRepo *repository.SeatRepository
}

func NewSeatService(seatRepo *repository.SeatRepository) *SeatService {
	return &SeatService{seatRepo: seatRepo}
}

func (s *SeatService) OccupySeat(seatID, userID string) error {
	return s.seatRepo.OccupySeat(seatID, userID)
}

func (s *SeatService) VacateSeat(seatID string) error {
	return s.seatRepo.VacateSeat(seatID)
}

func (s *SeatService) GetSeat(seatID string) (*domain.Seat, error) {
	return s.seatRepo.GetSeatByID(seatID)
}
