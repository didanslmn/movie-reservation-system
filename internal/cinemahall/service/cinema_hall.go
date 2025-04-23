package service

import (
	"context"
	"errors"

	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/repository"
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
	hall := model.CinemaHall{
		Name:     req.Name,
		Capacity: req.Capacity,
	}
	if err := s.repo.Create(ctx, &hall); err != nil {
		return nil, err
	}
	return mapper.ToCinemaHallResponse(&hall), nil
}

func (s *cinemaHallService) GetByID(ctx context.Context, id uint) (*response.CinemaHallResponse, error) {
	hall, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ToCinemaHallResponse(hall), nil
}

func (s *cinemaHallService) GetAll(ctx context.Context) ([]response.CinemaHallResponse, error) {
	halls, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.ToCinemaHallResponses(halls), nil
}

func (s *cinemaHallService) Update(ctx context.Context, id uint, req request.UpdateCinemaHallRequest) (*response.CinemaHallResponse, error) {
	hall, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("cinema hall not found")
	}

	hall.Name = req.Name
	hall.Capacity = req.Capacity

	if err := s.repo.Update(ctx, hall); err != nil {
		return nil, err
	}
	return mapper.ToCinemaHallResponse(hall), nil
}

func (s *cinemaHallService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
