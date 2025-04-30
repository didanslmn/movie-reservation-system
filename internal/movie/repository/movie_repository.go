package repository

import (
	"context"
	"errors"
	"fmt"

	genreModel "github.com/didanslmn/movie-reservation-system.git/internal/genre/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/model"
	"github.com/didanslmn/movie-reservation-system.git/utils"
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
	ExistsByID(ctx context.Context, id uint) (bool, error)
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
			utils.ErrorLogger.Printf("failed to validate genres: %v", err)
			return fmt.Errorf("failed to validate genres: %w", err)
		}
		if len(genres) != len(genreIDs) {
			return ErrNotFound
		}
		// Create movie
		if err := tx.Create(movie).Error; err != nil {
			utils.ErrorLogger.Printf("failed to create movie: %v", err)
			return fmt.Errorf("failed to create movie: %w", err)
		}
		// Associate genres
		if err := tx.Model(movie).Association("Genres").Append(genres); err != nil {
			utils.ErrorLogger.Printf("failed to associate genres: %v", err)
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
		utils.ErrorLogger.Printf("failed to get movie by id: %v", err)
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
		utils.ErrorLogger.Printf("failed to get all movies: %v", err)
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}
	return movies, nil
}

func (r *movieRepository) Update(ctx context.Context, movie *model.Movie) error {
	if err := r.db.WithContext(ctx).Save(movie).Error; err != nil {
		utils.ErrorLogger.Printf("failed to save movie (id: %d): %v", movie.ID, err)
		return fmt.Errorf("failed to save movie: %w", err)
	}
	return nil
}

func (r *movieRepository) UpdateGenres(ctx context.Context, movieID uint, genreIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// validasi genre
		var genres []genreModel.Genre
		if err := tx.Find(&genres, "id IN ?", genreIDs).Error; err != nil {
			utils.ErrorLogger.Printf("failed to validate genres for update: %v", err)
			return fmt.Errorf("fail to validate genre: %w", err)
		}
		if len(genres) != len(genreIDs) {
			return ErrNotFound
		}
		// Update associations
		movie := model.Movie{Model: gorm.Model{ID: movieID}}
		if err := tx.Model(&movie).Association("Genres").Replace(genres); err != nil {
			utils.ErrorLogger.Printf("failed to update movie genres: %v", err)
			return fmt.Errorf("failed to update genres: %w", err)
		}

		return nil
	})
}
func (r *movieRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).
		Select("Genres"). // Untuk cascade delete associations (bergantung pada konfigurasi model)
		Delete(&model.Movie{}, id)

	if res.Error != nil {
		utils.ErrorLogger.Printf("failed to delete movie (id: %d): %v", id, res.Error)
		return fmt.Errorf("failed to delete movie: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		utils.ErrorLogger.Printf("Delete movie failed - movie not Found (id: %d)", id)
		return errors.New("movie not found")
	}

	return nil
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
		utils.ErrorLogger.Printf("failed to get movies by genre: %v", err)
		return nil, fmt.Errorf("failed to get movies by genre: %w", err)
	}
	return movies, nil
}
func (r *movieRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Movie{}).Where("id = ?", id).Count(&count).Error; err != nil {
		utils.ErrorLogger.Printf("Failed to check existence of cinemahall ID %d: %v", id, err)
		return false, fmt.Errorf("failed to check cinemahall existence: %w", err)
	}
	return count > 0, nil
}
