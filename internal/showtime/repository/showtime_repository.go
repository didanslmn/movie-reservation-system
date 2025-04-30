package repository

import (
	"context"
	"fmt"

	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/model"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"gorm.io/gorm"
)

type ShowtimeRepository interface {
	Create(ctx context.Context, showtime *model.Showtime) error
	GetByID(ctx context.Context, id uint) (*model.Showtime, error)
	GetAll(ctx context.Context) ([]model.Showtime, error)
	Update(ctx context.Context, showtime *model.Showtime) error
	Delete(ctx context.Context, id uint) error
}

type showtimeRepository struct {
	db *gorm.DB
}

func NewShowtimeRepository(db *gorm.DB) ShowtimeRepository {
	return &showtimeRepository{db: db}
}

func (r *showtimeRepository) Create(ctx context.Context, showtime *model.Showtime) error {
	if err := r.db.WithContext(ctx).Create(showtime).Error; err != nil {
		utils.ErrorLogger.Printf("Failed to create showtime: %v", err)
		return fmt.Errorf("failed to create showtime: %w", err)
	}
	return nil
}

func (r *showtimeRepository) GetByID(ctx context.Context, id uint) (*model.Showtime, error) {
	var showtime model.Showtime
	err := r.db.WithContext(ctx).First(&showtime, id).Error
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get showtime by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get showtime by id: %w", err)
	}
	return &showtime, nil
}

func (r *showtimeRepository) GetAll(ctx context.Context) ([]model.Showtime, error) {
	var showtimes []model.Showtime
	err := r.db.WithContext(ctx).Find(&showtimes).Error
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get all showtimes: %v", err)
		return nil, fmt.Errorf("failed to get all showtime: %w", err)
	}
	return showtimes, nil
}

func (r *showtimeRepository) Update(ctx context.Context, showtime *model.Showtime) error {
	if err := r.db.WithContext(ctx).Save(showtime).Error; err != nil {
		utils.ErrorLogger.Printf("Failed to update showtime (ID: %d): %v", showtime.ID, err)
		return fmt.Errorf("failed to update showtime %w", err)
	}
	return nil
}

func (r *showtimeRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Showtime{}, id).Error; err != nil {
		utils.ErrorLogger.Printf("Failed to delete showtime (ID: %d): %v", id, err)
		return fmt.Errorf("failed to delete showtime: %w", err)
	}
	return nil
}
