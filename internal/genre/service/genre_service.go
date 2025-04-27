package service

import (
	"context"
	"fmt"

	"github.com/didanslmn/movie-reservation-system.git/internal/genre/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/repository"
	"github.com/didanslmn/movie-reservation-system.git/utils"
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
		if err == repository.ErrAlredyExists {
			utils.ErrorLogger.Printf("genre alredy exists : %s", req.Name)
			return nil, fmt.Errorf("genre alredy exists: %s", req.Name)
		}
		return nil, fmt.Errorf("failed to create genre: %w", err)
	}
	utils.InfoLogger.Printf("Created genre: %s", genre.Name)
	return mapper.ToGenreResponse(&genre), nil
}

func (s *genreService) GetGenre(ctx context.Context, id uint) (*response.Genre, error) {
	genre, err := s.repo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("Error fetching genre (id: %d): %v", id, err)
		return nil, fmt.Errorf("failed to get genre with ID %d: %w", id, err)
	}
	if genre == nil {
		utils.ErrorLogger.Printf("Genre not found (id: %d)", id)
		return nil, fmt.Errorf("genre with ID %d not found", id)
	}
	utils.InfoLogger.Printf("success get genre: %s", genre.Name)
	return mapper.ToGenreResponse(genre), nil
}

func (s *genreService) GetAllGenres(ctx context.Context) ([]response.Genre, error) {
	genres, err := s.repo.GetAll(ctx)
	if err != nil {
		utils.ErrorLogger.Printf("error fetching all genres: %v", err)
		return nil, fmt.Errorf("failed to fetch all genres: %w", err)
	}
	utils.InfoLogger.Println("Successfully fetched all genres")
	return mapper.ToGenreResponseList(genres), nil
}

func (s *genreService) UpdateGenre(ctx context.Context, id uint, req request.UpdateGenre) (*response.Genre, error) {
	existingGenre, err := s.repo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("genre not found (id: %d): %v", id, err)
		return nil, fmt.Errorf("genre not found (id: %d): %w", id, err)
	}

	existingGenre.Name = req.Name // Update field yang diperlukan

	if err := s.repo.Update(ctx, id, existingGenre); err != nil {
		utils.ErrorLogger.Printf("error updating genre (id: %d):%v", id, err)
		return nil, fmt.Errorf("failed to update genre: %w", err)
	}

	utils.InfoLogger.Println("Successfully update genre")
	return mapper.ToGenreResponse(existingGenre), nil
}

func (s *genreService) DeleteGenre(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		utils.ErrorLogger.Printf("error deleting genre (id: %d):%v", id, err)
		return fmt.Errorf("failedto delete genre: %w", err)
	}
	utils.InfoLogger.Printf("Successfully deleted genre (ID: %d)", id)
	return nil
}
