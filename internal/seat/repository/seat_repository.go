package repository

import (
	"context"
	"fmt"

	"github.com/didanslmn/movie-reservation-system.git/internal/seat/model"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"gorm.io/gorm"
)

type SeatRepository interface {
	Create(ctx context.Context, seat *model.Seat) error
	GetByID(ctx context.Context, id uint) (*model.Seat, error)
	GetByHallID(ctx context.Context, hallID uint) ([]model.Seat, error)
	Update(ctx context.Context, seat *model.Seat) error
	Delete(ctx context.Context, id uint) error
	// ExistsByID(ctx context.Context, id uint) (bool, error)
	GetByIDs(ctx context.Context, ids []uint) ([]model.Seat, error)
	IsSeatAvailable(seatID uint, showtimeID uint) (bool, error)
	UpdateStatus(ctx context.Context, seatID uint, status string) error
}

type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) Create(ctx context.Context, seat *model.Seat) error {
	if err := r.db.WithContext(ctx).Create(seat).Error; err != nil {
		utils.ErrorLogger.Printf("failed to create seat: %v", err)
		return fmt.Errorf("failed to create seat: %w", err)
	}
	return nil
}

func (r *seatRepository) GetByID(ctx context.Context, id uint) (*model.Seat, error) {
	var seat model.Seat
	if err := r.db.WithContext(ctx).First(&seat, id).Error; err != nil {
		utils.ErrorLogger.Printf("failed to get seat by id %d: %v", id, err)
		return nil, fmt.Errorf("failed to get seat by ID %w", err)
	}
	return &seat, nil
}

func (r *seatRepository) GetByHallID(ctx context.Context, hallID uint) ([]model.Seat, error) {
	var seats []model.Seat
	if err := r.db.WithContext(ctx).Where("cinema_hall_id = ?", hallID).Find(&seats).Error; err != nil {
		utils.ErrorLogger.Printf("failed to get seats by hall id %d: %v", hallID, err)
		return nil, fmt.Errorf("failed get by hall: %w", err)
	}
	return seats, nil
}

func (r *seatRepository) Update(ctx context.Context, seat *model.Seat) error {
	if err := r.db.WithContext(ctx).Save(seat).Error; err != nil {
		utils.ErrorLogger.Printf("failed to update seat id %d: %v", seat.ID, err)
		return err
	}
	return nil
}

func (r *seatRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Seat{}, id).Error; err != nil {
		utils.ErrorLogger.Printf("failed to delete seat id %d: %v", id, err)
		return err
	}
	return nil
}

//	func (r *seatRepository) ExistsByID(ctx context.Context, id uint) (bool, error) {
//		var count int64
//		if err := r.db.WithContext(ctx).Model(&model.Seat{}).Where("id = ?", id).Count(&count).Error; err != nil {
//			utils.ErrorLogger.Printf("Failed to check existence of seat ID %d: %v", id, err)
//			return false, fmt.Errorf("failed to check seat existence: %w", err)
//		}
//		return count > 0, nil
//	}
func (r *seatRepository) GetByIDs(ctx context.Context, ids []uint) ([]model.Seat, error) {
	var seats []model.Seat
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *seatRepository) IsSeatAvailable(seatID uint, showtimeID uint) (bool, error) {
	var count int64
	err := r.db.Table("reservations").
		Where("seat_id = ? AND showtime_id = ? AND is_active = TRUE", seatID, showtimeID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count == 0, nil
}
func (r *seatRepository) UpdateStatus(ctx context.Context, seatID uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.Seat{}).
		Where("id = ?", seatID).
		Update("status", status).Error
}
