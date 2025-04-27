package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/didanslmn/movie-reservation-system.git/internal/genre/model"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"gorm.io/gorm"
)

type GenreRepository interface {
	Create(ctx context.Context, genre *model.Genre) error
	GetByID(ctx context.Context, id uint) (*model.Genre, error)
	GetAll(ctx context.Context) ([]model.Genre, error)
	Update(ctx context.Context, id uint, genre *model.Genre) error
	Delete(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByIDs(ctx context.Context, ids []uint) (bool, error)
}

type genreRepository struct {
	db *gorm.DB
}

// NewRepositoryGenre mengembalikan implementasi GenreRepository
func NewRepositoryGenre(db *gorm.DB) GenreRepository {
	return &genreRepository{db: db}
}

var (
	ErrNotFound     = errors.New("not found")
	ErrAlredyExists = errors.New("genre alredy exists")
)

func (r *genreRepository) Create(ctx context.Context, genre *model.Genre) error {
	// cek genre untuk nama yang sama
	if exists, _ := r.ExistsByName(ctx, genre.Name); exists {
		return ErrAlredyExists
	}
	// simpan ke db
	if err := r.db.WithContext(ctx).Create(genre).Error; err != nil { //.Create(&genre): GORM akan menjalankan query INSERT ke dalam tabel sesuai dengan struct genre (biasanya ke tabel genres).
		utils.ErrorLogger.Printf("Create genre failed (genre: %s):%v", genre.Name, err)
		return fmt.Errorf("failed to create genre: %w", err)
	}
	return nil
}
func (r *genreRepository) GetByID(ctx context.Context, id uint) (*model.Genre, error) {
	var genre model.Genre
	if err := r.db.WithContext(ctx).First(&genre, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		utils.ErrorLogger.Printf("Get genre bt id failed (id:%d):%v", id, err)
		return nil, fmt.Errorf("failed to get genre by ID %w", err)
	}
	return &genre, nil
}
func (r *genreRepository) GetAll(ctx context.Context) ([]model.Genre, error) {
	var genres []*model.Genre
	if err := r.db.WithContext(ctx).Find(&genres).Error; err != nil {
		utils.ErrorLogger.Printf("Failed to get all genres: %v", err)
		return nil, fmt.Errorf("failed to get genres: %w", err)
	}

	result := make([]model.Genre, len(genres))
	for i, genre := range genres {
		result[i] = *genre
	}
	return result, nil
}
func (r *genreRepository) Update(ctx context.Context, id uint, input *model.Genre) error {
	var genre *model.Genre
	if err := r.db.WithContext(ctx).First(&genre, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorLogger.Printf("update genre failed - not found (id:%d)", id)
			return ErrNotFound
		}
		utils.ErrorLogger.Printf("update genre failed (id:%d):%v", id, err)
		return fmt.Errorf("failed to find genre by ID: %w", err)
	}

	if genre.Name != input.Name {
		exists, err := r.ExistsByName(ctx, input.Name)
		if err != nil {
			utils.ErrorLogger.Printf("Check name exists failed (name: %s): %v", input.Name, err)
			return fmt.Errorf("failed to check existing genre name: %w", err)
		}
		if exists {
			utils.ErrorLogger.Printf("Update genre failed - name already exists (name: %s)", input.Name)
			return ErrAlredyExists
		}
		genre.Name = input.Name
	}

	// Simpan perubahan
	if err := r.db.WithContext(ctx).Save(&genre).Error; err != nil {
		utils.ErrorLogger.Printf("Save updated genre failed (id: %d): %v", id, err)
		return fmt.Errorf("failed to save updated genre: %w", err)

	}

	return nil
}

func (r *genreRepository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Delete(&model.Genre{}, id)
	if res.Error != nil {
		utils.ErrorLogger.Printf("Delete genre failed (id: %d):%v", id, res.Error)
		return fmt.Errorf("failed to delete genre: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		utils.ErrorLogger.Printf("Delete genre failed - genre not Found (id: %d)", id)
		return ErrNotFound
	}
	return nil
}
func (r *genreRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Genre{}).
		Where("name = ?", name).
		Count(&count).
		Error
	if err != nil {
		utils.ErrorLogger.Printf("check if genre exist by name failed (name: %s):%v", name, err)
		return false, fmt.Errorf("failed to check if genre exists by name: %w", err)
	}
	return count > 0, err
}
func (r *genreRepository) ExistsByIDs(ctx context.Context, ids []uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Genre{}).
		Where("id IN ?", ids).
		Count(&count).Error; err != nil {
		utils.ErrorLogger.Printf("check if genre exist by IDs failed (id :%d):%v", ids, err)
		return false, fmt.Errorf("failed to check if genre exists by IDs: %w", err)
	}
	return count == int64(len(ids)), nil
}
