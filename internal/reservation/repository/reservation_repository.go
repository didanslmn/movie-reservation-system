package repository

import (
	"context"
	"fmt"

	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/model"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation *model.Reservation) error
	GetByID(ctx context.Context, id uint) (*model.Reservation, error)
	GetAll(ctx context.Context) ([]model.Reservation, error)
	Delete(ctx context.Context, id uint) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(ctx context.Context, reservation *model.Reservation) error {
	if err := r.db.WithContext(ctx).Create(reservation).Error; err != nil {
		utils.ErrorLogger.Printf("Error to create reservation: %v", err)
		return fmt.Errorf("failed to create reservation: %w", err)
	}
	return nil
}

func (r *reservationRepository) GetByID(ctx context.Context, id uint) (*model.Reservation, error) {
	var reservation model.Reservation
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Showtime").
		Preload("Seat").
		First(&reservation, id).Error
	if err != nil {
		utils.ErrorLogger.Printf("Error to get reservation by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get reservation by id: %w", err)
	}
	return &reservation, nil
}

func (r *reservationRepository) GetAll(ctx context.Context) ([]model.Reservation, error) {
	var reservations []model.Reservation
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Showtime").
		Preload("Seat").
		Find(&reservations).Error
	if err != nil {
		utils.ErrorLogger.Printf("Error to get all reservations: %v", err)
		return nil, fmt.Errorf("failed to get all reservations: %w", err)
	}
	return reservations, nil
}

func (r *reservationRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Reservation{}, id).Error; err != nil {
		utils.ErrorLogger.Printf("Error to delete reservation (ID: %d): %v", id, err)
		return fmt.Errorf("failed to delete reservation: %w", err)
	}
	return nil
}
