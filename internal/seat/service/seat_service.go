package service

import (
	"context"
	"fmt"

	"github.com/didanslmn/movie-reservation-system.git/internal/seat/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/repository"
	"github.com/didanslmn/movie-reservation-system.git/utils"
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
	seat := mapper.ToSeatModel(&req)

	if err := s.repo.Create(ctx, seat); err != nil {
		utils.ErrorLogger.Printf("Error Creating create seat: %v", err)
		return nil, fmt.Errorf("failed to create seat: %w", err)
	}
	utils.InfoLogger.Printf("Success creat seat: %s", req.SeatNumber)
	return mapper.ToSeatResponse(seat), nil
}

func (s *seatService) GetSeat(ctx context.Context, id uint) (*response.SeatResponse, error) {
	seat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("failed to get seat id %d: %v", id, err)
		return nil, fmt.Errorf("failed to get seat by ID %d:%w", id, err)
	}
	utils.InfoLogger.Printf("Success fetch seat (ID: %d)", id)
	return mapper.ToSeatResponse(seat), nil
}

func (s *seatService) GetByHallID(ctx context.Context, hallID uint) ([]response.SeatResponse, error) {
	seats, err := s.repo.GetByHallID(ctx, hallID)
	if err != nil {
		utils.ErrorLogger.Printf("failed to get seats by hall id %d: %v", hallID, err)
		return nil, fmt.Errorf("failed to get seat by hall id: %w", err)
	}
	utils.InfoLogger.Printf("Success fetch seat by hall (ID: %d)", hallID)
	return mapper.ToSeatResponseList(seats), nil
}

func (s *seatService) UpdateSeat(ctx context.Context, id uint, req request.UpdateSeatRequest) (*response.SeatResponse, error) {
	seat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("failed to find seat id %d for update: %v", id, err)
		return nil, fmt.Errorf("failed to update seat: %w", err)
	}

	seat.SeatNumber = req.SeatNumber
	seat.Row = req.Row
	seat.Status = req.Status

	if err := s.repo.Update(ctx, seat); err != nil {
		utils.ErrorLogger.Printf("Error to updating seat id %d: %v", id, err)
		return nil, fmt.Errorf("failed to update seat: %w", err)
	}
	utils.InfoLogger.Println("Successfully update seats")
	return mapper.ToSeatResponse(seat), nil
}

func (s *seatService) DeleteSeat(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		utils.ErrorLogger.Printf("Error to delete seat id %d: %v", id, err)
		return fmt.Errorf("failed to delete seat :%w", err)
	}
	utils.InfoLogger.Printf("Successfully deleted seat (ID: %d)", id)
	return nil
}
