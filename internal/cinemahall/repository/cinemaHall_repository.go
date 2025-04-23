package repository

import (
	"context"

	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/model"
	"gorm.io/gorm"
)

type CinemaHallRepository interface {
	Create(ctx context.Context, hall *model.CinemaHall) error
	GetByID(ctx context.Context, id uint) (*model.CinemaHall, error)
	GetAll(ctx context.Context) ([]model.CinemaHall, error)
	Update(ctx context.Context, hall *model.CinemaHall) error
	Delete(ctx context.Context, id uint) error
}

type cinemaHallRepository struct {
	db *gorm.DB
}

func NewCinemaHallRepository(db *gorm.DB) CinemaHallRepository {
	return &cinemaHallRepository{db: db}
}

func (r *cinemaHallRepository) Create(ctx context.Context, hall *model.CinemaHall) error {
	return r.db.WithContext(ctx).Create(hall).Error
}

func (r *cinemaHallRepository) GetByID(ctx context.Context, id uint) (*model.CinemaHall, error) {
	var hall *model.CinemaHall
	err := r.db.WithContext(ctx).Preload("Seats").Find(&hall, id).Error
	return hall, err
}
func (r *cinemaHallRepository) GetAll(ctx context.Context) ([]model.CinemaHall, error) {
	var halls []model.CinemaHall
	err := r.db.WithContext(ctx).Preload("Seats").Find(&halls).Error
	return halls, err
}
func (r *cinemaHallRepository) Update(ctx context.Context, hall *model.CinemaHall) error {
	return r.db.WithContext(ctx).Save(hall).Error
}
func (r *cinemaHallRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.CinemaHall{}, "id = ?", id).Error
}
