package repository

import (
	"context"
	"errors"

	"github.com/didanslmn/movie-reservation-system.git/internal/genre/model"
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
	ErrNotFound = errors.New("not found")
)

func (r *genreRepository) Create(ctx context.Context, genre *model.Genre) error {
	// cek genre untuk nama yang sama
	if exists, _ := r.ExistsByName(ctx, genre.Name); exists {
		return errors.New("genre alredy exists")
	}
	// simpan ke db
	if err := r.db.WithContext(ctx).Create(genre).Error; err != nil { //.Create(&genre): GORM akan menjalankan query INSERT ke dalam tabel sesuai dengan struct genre (biasanya ke tabel genres).
		return err
	}
	return nil
}
func (r *genreRepository) GetByID(ctx context.Context, id uint) (*model.Genre, error) {
	var genre model.Genre
	if err := r.db.WithContext(ctx).First(&genre, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("genre not found")
		}
		return nil, err
	}
	return &genre, nil
}
func (r *genreRepository) GetAll(ctx context.Context) ([]model.Genre, error) {
	var genres []*model.Genre
	if err := r.db.WithContext(ctx).Find(&genres).Error; err != nil {
		return nil, err
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
			return errors.New("genre not found")
		}
		return err
	}

	if genre.Name != input.Name {
		exists, err := r.ExistsByName(ctx, input.Name)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("new genre name already exists")
		}
		genre.Name = input.Name
	}

	// Simpan perubahan
	if err := r.db.WithContext(ctx).Save(&genre).Error; err != nil {
		return err
	}

	return nil
}

func (r *genreRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Genre{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("genre not found")
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
	return count > 0, err
}
func (r *genreRepository) ExistsByIDs(ctx context.Context, ids []uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Genre{}).
		Where("id IN ?", ids).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count == int64(len(ids)), nil
}
