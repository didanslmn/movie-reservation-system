package service

import (
	"context"

	"github.com/didanslmn/movie-reservation-system.git/internal/genre/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/repository"
)

type GenreService interface {
	CreateGenre(ctx context.Context, req request.CreateGenre) (*response.Genre, error)
	GetGenre(ctx context.Context, id uint) (*response.Genre, error)
	GetAllGenres(ctx context.Context) ([]response.Genre, error)
	UpdateGenre(ctx context.Context, id uint, req request.UpdateGenre) (*response.Genre, error)
	DeleteGenre(ctx context.Context, id uint) error
}

type genreService struct {
	repo repository.GenreRepository
}

func NewGenreService(repo repository.GenreRepository) GenreService {
	return &genreService{repo: repo}
}

func (s *genreService) CreateGenre(ctx context.Context, req request.CreateGenre) (*response.Genre, error) {
	genre := model.Genre{
		Name: req.Name,
	}

	if err := s.repo.Create(ctx, &genre); err != nil {
		return nil, err
	}

	return mapper.ToGenreResponse(&genre), nil
}

func (s *genreService) GetGenre(ctx context.Context, id uint) (*response.Genre, error) {
	genre, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ToGenreResponse(genre), nil
}

func (s *genreService) GetAllGenres(ctx context.Context) ([]response.Genre, error) {
	genres, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.ToGenreResponseList(genres), nil
}

func (s *genreService) UpdateGenre(ctx context.Context, id uint, req request.UpdateGenre) (*response.Genre, error) {
	genre := model.Genre{
		Name: req.Name,
	}

	if err := s.repo.Update(ctx, id, &genre); err != nil {
		return nil, err
	}

	updatedGenre, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.ToGenreResponse(updatedGenre), nil
}

func (s *genreService) DeleteGenre(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
