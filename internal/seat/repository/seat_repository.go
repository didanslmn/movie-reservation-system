package repository

import (
	"context"

	"github.com/didanslmn/movie-reservation-system.git/internal/seat/model"
	"gorm.io/gorm"
)

type SeatRepository interface {
	Create(ctx context.Context, seat *model.Seat) error
	GetByID(ctx context.Context, id uint) (*model.Seat, error)
	GetByHallID(ctx context.Context, hallID uint) ([]model.Seat, error)
	Update(ctx context.Context, seat *model.Seat) error
	Delete(ctx context.Context, id uint) error
}
type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) Create(ctx context.Context, seat *model.Seat) error {
	return r.db.WithContext(ctx).Create(seat).Error
}

func (r *seatRepository) GetByID(ctx context.Context, id uint) (*model.Seat, error) {
	var seat model.Seat
	err := r.db.WithContext(ctx).First(&seat, id).Error
	return &seat, err
}

func (r *seatRepository) GetByHallID(ctx context.Context, hallID uint) ([]model.Seat, error) {
	var seats []model.Seat
	err := r.db.WithContext(ctx).Where("cinema_hall_id = ?", hallID).Find(&seats).Error
	return seats, err
}

func (r *seatRepository) Update(ctx context.Context, seat *model.Seat) error {
	return r.db.WithContext(ctx).Save(seat).Error
}

func (r *seatRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Seat{}, id).Error
}
