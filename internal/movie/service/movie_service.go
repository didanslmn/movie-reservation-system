package service

import (
	"context"
	"errors"

	genreRepository "github.com/didanslmn/movie-reservation-system.git/internal/genre/repository"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/repository"
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
	// Validate genre IDs
	exist, err := s.genreRepo.ExistsByIDs(ctx, req.GenreIDs)
	if err != nil || !exist {
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
		return nil, err
	}

	return s.GetMovie(ctx, movie.ID)
}

func (s *movieService) GetMovie(ctx context.Context, id uint) (*response.Movie, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapper.ToMovieResponse(movie), nil
}

func (s *movieService) GetAllMovies(ctx context.Context) ([]response.Movie, error) {
	movies, err := s.movieRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]response.Movie, len(movies))
	for i, m := range movies {
		responses[i] = *mapper.ToMovieResponse(&m)
	}

	return responses, nil
}

func (s *movieService) UpdateMovie(ctx context.Context, id uint, req request.UpdateMovie) (*response.Movie, error) {
	movie, err := s.movieRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// cek field update
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

	if err := s.movieRepo.Update(ctx, movie); err != nil {
		return nil, err
	}
	// Update genres
	if req.GenreIDs != nil {
		if err := s.movieRepo.UpdateGenres(ctx, id, req.GenreIDs); err != nil {
			return nil, err
		}
	}

	return s.GetMovie(ctx, id)
}

func (s *movieService) DeleteMovie(ctx context.Context, id uint) error {
	return s.movieRepo.Delete(ctx, id)
}

func (s *movieService) GetMoviesByGenre(ctx context.Context, genreID uint) ([]response.Movie, error) {
	movies, err := s.movieRepo.GetByGenre(ctx, genreID)
	if err != nil {
		return nil, err
	}

	responses := make([]response.Movie, len(movies))
	for i, m := range movies {
		responses[i] = *mapper.ToMovieResponse(&m)
	}

	return responses, nil
}
