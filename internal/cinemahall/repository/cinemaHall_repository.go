package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/model"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"gorm.io/gorm"
)

type CinemaHallRepository interface {
	Create(ctx context.Context, hall *model.CinemaHall) error
	GetByID(ctx context.Context, id uint) (*model.CinemaHall, error)
	GetAll(ctx context.Context) ([]model.CinemaHall, error)
	Update(ctx context.Context, hall *model.CinemaHall, id uint) error
	Delete(ctx context.Context, id uint) error
}

type cinemaHallRepository struct {
	db *gorm.DB
}

func NewCinemaHallRepository(db *gorm.DB) CinemaHallRepository {
	return &cinemaHallRepository{db: db}
}

func (r *cinemaHallRepository) Create(ctx context.Context, hall *model.CinemaHall) error {
	// Cek jika cinema hall sudah ada dengan nama yang sama
	if err := r.db.WithContext(ctx).Where("name = ?", hall.Name).First(&model.CinemaHall{}).Error; err == nil {
		utils.ErrorLogger.Printf("Cinema hall with this name already exists")
		return fmt.Errorf("genre alredy exists: %w", err)
	}
	if err := r.db.WithContext(ctx).Create(hall).Error; err != nil {
		utils.ErrorLogger.Printf("Failed to create cinema hall (hall: %s):%v", hall.Name, err)
		return fmt.Errorf("failed to create cinema hall: %w", err)
	}

	return nil
}

func (r *cinemaHallRepository) GetByID(ctx context.Context, id uint) (*model.CinemaHall, error) {
	var hall model.CinemaHall
	// Mengambil data cinema hall berdasarkan ID
	if err := r.db.WithContext(ctx).Preload("Seat").First(&hall, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cinema hall not found")
		}
		utils.ErrorLogger.Printf("Get genre by id failed (id: %d):%v", id, err)
		return nil, fmt.Errorf("failed get genre by id: %w", err)
	}

	return &hall, nil
}

func (r *cinemaHallRepository) GetAll(ctx context.Context) ([]model.CinemaHall, error) {
	var halls []model.CinemaHall
	// Mengambil semua cinema hall
	if err := r.db.WithContext(ctx).Preload("Seat").Find(&halls).Error; err != nil {
		utils.ErrorLogger.Printf("Failed to get all genre:%v ", err)
		return nil, fmt.Errorf("failed to get genres: %w", err)
	}
	return halls, nil
}

func (r *cinemaHallRepository) Update(ctx context.Context, hall *model.CinemaHall, id uint) error {
	// Memperbarui cinema hall
	if err := r.db.WithContext(ctx).Save(hall).Error; err != nil {
		utils.ErrorLogger.Printf("Failed to update cinema hall (id: %d):%v ", id, err)
		return err
	}

	return nil
}

func (r *cinemaHallRepository) Delete(ctx context.Context, id uint) error {

	res := r.db.WithContext(ctx).Delete(&model.CinemaHall{}, id)
	if res.Error != nil {
		utils.ErrorLogger.Printf("Failed to delete cinema hall (id: %d):%v ", id, res.Error)
		return fmt.Errorf("failed to delete cinema hall:%w", res.Error)
	}
	if res.RowsAffected == 0 {
		utils.ErrorLogger.Printf("Delete genre failed - cinema hall not Found (id: %d)", id)
		return errors.New("cinema hall not found")
	}
	return nil
}
