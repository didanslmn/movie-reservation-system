package service

import (
	"context"
	"fmt"

	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/repository"
	"github.com/didanslmn/movie-reservation-system.git/utils"
)

type CinemaHallService interface {
	Create(ctx context.Context, req request.CreateCinemaHallRequest) (*response.CinemaHallResponse, error)
	GetByID(ctx context.Context, id uint) (*response.CinemaHallResponse, error)
	GetAll(ctx context.Context) ([]response.CinemaHallResponse, error)
	Update(ctx context.Context, id uint, req request.UpdateCinemaHallRequest) (*response.CinemaHallResponse, error)
	Delete(ctx context.Context, id uint) error
}

type cinemaHallService struct {
	repo repository.CinemaHallRepository
}

func NewCinemaHallService(repo repository.CinemaHallRepository) CinemaHallService {
	return &cinemaHallService{repo: repo}
}

func (s *cinemaHallService) Create(ctx context.Context, req request.CreateCinemaHallRequest) (*response.CinemaHallResponse, error) {
	hall := mapper.ToCinemaHallModel(&req)

	if err := s.repo.Create(ctx, hall); err != nil {
		utils.ErrorLogger.Printf("cinema hall alredy exists : %s", req.Name)
		return nil, fmt.Errorf("failed to create cinema hall: %w", err)
	}
	utils.InfoLogger.Printf("Created cinema hall: %s", hall.Name)
	return mapper.ToCinemaHallResponse(hall), nil
}

func (s *cinemaHallService) GetByID(ctx context.Context, id uint) (*response.CinemaHallResponse, error) {
	hall, err := s.repo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("Error fetching cinemahall (id: %d): %v", id, err)
		return nil, fmt.Errorf("failed to get cinema hall by ID %d: %w", id, err)
	}
	if hall == nil {
		utils.ErrorLogger.Printf("cinema hall not found (id: %d)", id)
		return nil, fmt.Errorf("cinema hall with ID %d not found", id)
	}
	utils.InfoLogger.Printf("success get genre: %s", hall.Name)
	return mapper.ToCinemaHallResponse(hall), nil
}

func (s *cinemaHallService) GetAll(ctx context.Context) ([]response.CinemaHallResponse, error) {
	halls, err := s.repo.GetAll(ctx)
	if err != nil {
		utils.ErrorLogger.Printf("error fetching all cinema: %v", err)
		return nil, fmt.Errorf("failed to get all cinema halls: %w", err)
	}
	utils.InfoLogger.Println("Successfully fetched all cinemas")
	return mapper.ToCinemaHallResponses(halls), nil
}

func (s *cinemaHallService) Update(ctx context.Context, id uint, req request.UpdateCinemaHallRequest) (*response.CinemaHallResponse, error) {
	hall, err := s.repo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("cinema hall not found (id: %d): %v", id, err)
		return nil, fmt.Errorf("cinema hall not found (id: %d): %w", id, err)
	}

	mapper.UpdateCinemaHallModel(hall, &req)

	if err := s.repo.Update(ctx, hall, id); err != nil {
		utils.ErrorLogger.Printf("error updating cinema hall (id: %d):%v", id, err)
		return nil, fmt.Errorf("failed to update cinema hall (id: %d): %w", id, err)
	}
	utils.InfoLogger.Println("Successfully update cinema hall")
	return mapper.ToCinemaHallResponse(hall), nil
}

func (s *cinemaHallService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		utils.ErrorLogger.Printf("error deleting cinema hall (id: %d):%v", id, err)
		return fmt.Errorf("failed to delete cinema hall (id: %d): %w", id, err)
	}
	utils.InfoLogger.Printf("Successfully deleted cinema hall (ID: %d)", id)
	return nil
}
