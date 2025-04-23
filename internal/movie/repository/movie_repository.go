package repository

import (
	"context"
	"errors"
	"fmt"

	genreModel "github.com/didanslmn/movie-reservation-system.git/internal/genre/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/model"
	"gorm.io/gorm"
)

type MovieRepository interface {
	Create(ctx context.Context, movie *model.Movie, genreIDs []uint) error
	GetByID(ctx context.Context, id uint) (*model.Movie, error)
	GetAll(ctx context.Context) ([]model.Movie, error)
	Update(ctx context.Context, movie *model.Movie) error
	UpdateGenres(ctx context.Context, movieID uint, genreIDs []uint) error
	Delete(ctx context.Context, id uint) error
	GetByGenre(ctx context.Context, genreID uint) ([]model.Movie, error)
}

type movieRepository struct {
	db *gorm.DB
}

func NewRepositoryMovie(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

var (
	ErrNotFound = errors.New("not found")
)

func (r *movieRepository) Create(ctx context.Context, movie *model.Movie, genreIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Validasi genre
		var genres []genreModel.Genre
		if err := tx.Find(&genres, "id IN ?", genreIDs).Error; err != nil {
			return fmt.Errorf("failed to validate genres: %w", err)
		}
		if len(genres) != len(genreIDs) {
			return errors.New("some genres not found")
		}
		// Create movie
		if err := tx.Create(movie).Error; err != nil {
			return fmt.Errorf("failed to create movie: %w", err)
		}
		// Associate genres
		if err := tx.Model(movie).Association("Genres").Append(genres); err != nil {
			return fmt.Errorf("failed to associate genres: %w", err)
		}

		return nil
	})
}

func (r *movieRepository) GetByID(ctx context.Context, id uint) (*model.Movie, error) {
	var movie model.Movie
	err := r.db.WithContext(ctx).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).First(&movie, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("movie with id %d: %w", id, ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get movie: %w", err)
	}
	return &movie, err
}

func (r *movieRepository) GetAll(ctx context.Context) ([]model.Movie, error) {
	var movies []model.Movie
	err := r.db.WithContext(ctx).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).Find(&movies).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}
	return movies, nil
}

func (r *movieRepository) Update(ctx context.Context, movie *model.Movie) error {
	return r.db.WithContext(ctx).
		Model(movie).
		Updates(map[string]any{
			"title":        movie.Title,
			"description":  movie.Description,
			"duration":     movie.Duration,
			"release_date": movie.ReleaseDate,
			"image_url":    movie.ImageURL,
			"rating":       movie.Rating,
		}).Error
}

func (r *movieRepository) UpdateGenres(ctx context.Context, movieID uint, genreIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// validasi genre
		var genres []genreModel.Genre
		if err := tx.Find(&genres, "id IN ?", genreIDs).Error; err != nil {
			return fmt.Errorf("fail to validate genre: %w", err)
		}
		if len(genres) != len(genreIDs) {
			return errors.New("some genres not found")
		}
		// Update associations
		movie := model.Movie{Model: gorm.Model{ID: movieID}}
		if err := tx.Model(&movie).Association("Genres").Replace(genres); err != nil {
			return fmt.Errorf("failed to update genres: %w", err)
		}

		return nil
	})
}
func (r *movieRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Select("Genres"). // Untuk cascade delete associations
		Delete(&model.Movie{}, id).
		Error
}

func (r *movieRepository) GetByGenre(ctx context.Context, genreID uint) ([]model.Movie, error) {
	var movies []model.Movie
	err := r.db.WithContext(ctx).
		Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").
		Where("movie_genres.genre_id = ?", genreID).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Find(&movies).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to get movies by genre: %w", err)
	}
	return movies, nil
}
