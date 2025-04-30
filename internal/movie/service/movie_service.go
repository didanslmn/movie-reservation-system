package service

import (
	"context"
	"errors"
	"fmt"

	genreRepository "github.com/didanslmn/movie-reservation-system.git/internal/genre/repository"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/repository"
	"github.com/didanslmn/movie-reservation-system.git/utils"
)

type MovieService interface {
	CreateMovie(ctx context.Context, req request.CreateMovie) (*response.Movie, error)
	GetMovie(ctx context.Context, id uint) (*response.Movie, error)
	GetAllMovies(ctx context.Context) ([]response.Movie, error)
	UpdateMovie(ctx context.Context, id uint, req request.UpdateMovie) (*response.Movie, error)
	DeleteMovie(ctx context.Context, id uint) error
	GetMoviesByGenre(ctx context.Context, genreID uint) ([]response.Movie, error)
}

type movieService struct {
	movieRepo repository.MovieRepository
	genreRepo genreRepository.GenreRepository
}

func NewMovieService(movieR repository.MovieRepository, genreR genreRepository.GenreRepository) MovieService {
	return &movieService{
		movieRepo: movieR,
		genreRepo: genreR,
	}
}
func (s *movieService) CreateMovie(ctx context.Context, req request.CreateMovie) (*response.Movie, error) {
	exist, err := s.genreRepo.ExistsByIDs(ctx, req.GenreIDs)
	if err != nil || !exist {
		utils.ErrorLogger.Printf("Invalid genre IDs: %v", req.GenreIDs)
		return nil, errors.New("invalid genre IDs")
	}

	movie := model.Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		ReleaseDate: req.ReleaseDate,
		ImageURL:    req.ImageURL,
		Rating:      req.Rating,
	}

	if err := s.movieRepo.Create(ctx, &movie, req.GenreIDs); err != nil {
		utils.ErrorLogger.Printf("Error creating movie: %v", err)
		return nil, fmt.Errorf("failed to create movie: %w", err)
	}
	utils.InfoLogger.Printf("Successfully created movie: %s (ID: %d)", movie.Title, movie.ID)
	return s.GetMovie(ctx, movie.ID)
}

func (s *movieService) GetMovie(ctx context.Context, id uint) (*response.Movie, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("Error fetching movie (id: %d): %v", id, err)
		return nil, fmt.Errorf("failed to get movie by ID %d: %w", id, err)
	}
	if movie == nil {
		utils.ErrorLogger.Printf("Movie not found (id: %d): %v", id, err)
		return nil, fmt.Errorf("movie not found: %w", err)
	}
	utils.InfoLogger.Printf("Successfully fetched movie (ID: %d)", id)
	return mapper.ToMovieResponse(movie), nil
}

func (s *movieService) GetAllMovies(ctx context.Context) ([]response.Movie, error) {
	movies, err := s.movieRepo.GetAll(ctx)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to fetch all movies: %v", err)
		return nil, fmt.Errorf("failed to get all movies")
	}
	utils.InfoLogger.Println("Successfully fetched all movies")
	responses := make([]response.Movie, len(movies))
	for i, m := range movies {
		responses[i] = *mapper.ToMovieResponse(&m)
	}

	return responses, nil
}

func (s *movieService) UpdateMovie(ctx context.Context, id uint, req request.UpdateMovie) (*response.Movie, error) {
	// Cari movie berdasarkan ID
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		utils.ErrorLogger.Printf("movie not found (id: %d): %v", id, err)
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	// Update field-field yang diinputkan
	if req.Title != "" {
		movie.Title = req.Title
	}
	if req.Description != "" {
		movie.Description = req.Description
	}
	if req.Duration > 0 {
		movie.Duration = req.Duration
	}
	if !req.ReleaseDate.IsZero() {
		movie.ReleaseDate = req.ReleaseDate
	}
	if req.ImageURL != "" {
		movie.ImageURL = req.ImageURL
	}
	if req.Rating > 0 {
		movie.Rating = req.Rating
	}

	// Update movie ke database
	if err := s.movieRepo.Update(ctx, movie); err != nil {
		utils.ErrorLogger.Printf("failed to update movie (id: %d): %v", id, err)
		return nil, fmt.Errorf("failed to update movie: %w", err)
	}

	utils.InfoLogger.Printf("successfully updated movie: %s (id: %d)", movie.Title, movie.ID)

	// Update genre jika ada perubahan genre
	if req.GenreIDs != nil {
		exist, err := s.genreRepo.ExistsByIDs(ctx, req.GenreIDs)
		if err != nil {
			utils.ErrorLogger.Printf("failed to check genre existence: %v", err)
			return nil, fmt.Errorf("failed to validate genres: %w", err)
		}
		if !exist {
			utils.ErrorLogger.Printf("invalid genre IDs: %v", req.GenreIDs)
			return nil, fmt.Errorf("invalid genre IDs")
		}

		if err := s.movieRepo.UpdateGenres(ctx, id, req.GenreIDs); err != nil {
			utils.ErrorLogger.Printf("failed to update genres for movie (id: %d): %v", id, err)
			return nil, fmt.Errorf("failed to update genres: %w", err)
		}
	}

	// Return movie terbaru
	return s.GetMovie(ctx, id)
}

func (s *movieService) DeleteMovie(ctx context.Context, id uint) error {
	if err := s.movieRepo.Delete(ctx, id); err != nil {
		utils.ErrorLogger.Printf("Error deleting movie (id: %d): %v", id, err)
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	utils.InfoLogger.Printf("Successfully deleted movie (ID: %d)", id)

	return nil
}

func (s *movieService) GetMoviesByGenre(ctx context.Context, genreID uint) ([]response.Movie, error) {
	movies, err := s.movieRepo.GetByGenre(ctx, genreID)
	if err != nil {
		utils.ErrorLogger.Printf("Error fetching movies by genre (genre ID: %d): %v", genreID, err)
		return nil, fmt.Errorf("failed to get movies by genre: %w", err)
	}

	utils.InfoLogger.Printf("Successfully fetched movies for genre (ID: %d)", genreID)

	responses := make([]response.Movie, len(movies))
	for i, m := range movies {
		responses[i] = *mapper.ToMovieResponse(&m)
	}

	return responses, nil
}
