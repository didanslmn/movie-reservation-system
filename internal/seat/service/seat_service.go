package service

import (
	"context"

	"github.com/didanslmn/movie-reservation-system.git/internal/seat/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/repository"
)

type SeatService interface {
	CreateSeat(ctx context.Context, req request.CreateSeatRequest) (*response.SeatResponse, error)
	GetSeat(ctx context.Context, id uint) (*response.SeatResponse, error)
	GetByHallID(ctx context.Context, hallID uint) ([]response.SeatResponse, error)
	UpdateSeat(ctx context.Context, id uint, req request.UpdateSeatRequest) (*response.SeatResponse, error)
	DeleteSeat(ctx context.Context, id uint) error
}

type seatService struct {
	repo repository.SeatRepository
}

func NewSeatService(repo repository.SeatRepository) SeatService {
	return &seatService{repo: repo}
}

func (s *seatService) CreateSeat(ctx context.Context, req request.CreateSeatRequest) (*response.SeatResponse, error) {
	seat := model.Seat{
		CinemaHallID: req.CinemaHallID,
		SeatNumber:   req.SeatNumber,
		Row:          req.Row,
	}
	if err := s.repo.Create(ctx, &seat); err != nil {
		return nil, err
	}
	return mapper.ToSeatResponse(&seat), nil
}
func (s *seatService) GetSeat(ctx context.Context, id uint) (*response.SeatResponse, error) {
	seat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ToSeatResponse(seat), nil
}

func (s *seatService) GetByHallID(ctx context.Context, hallID uint) ([]response.SeatResponse, error) {
	seats, err := s.repo.GetByHallID(ctx, hallID)
	if err != nil {
		return nil, err
	}
	return mapper.ToSeatResponseList(seats), nil
}

func (s *seatService) UpdateSeat(ctx context.Context, id uint, req request.UpdateSeatRequest) (*response.SeatResponse, error) {
	seat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	seat.SeatNumber = req.SeatNumber
	seat.Row = req.Row
	seat.Status = req.Status

	if err := s.repo.Update(ctx, seat); err != nil {
		return nil, err
	}
	return mapper.ToSeatResponse(seat), nil
}

func (s *seatService) DeleteSeat(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
