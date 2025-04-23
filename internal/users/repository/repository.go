package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/didanslmn/movie-reservation-system.git/internal/pkg/logger"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id uint) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, id uint, newPassword string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		logger.Log.Error("failed to create user", zap.Error(err), zap.String("email", user.Email))
		return fmt.Errorf("failet to create user: %w", err)
	}
	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Log.Error("failed to get user by email", zap.Error(err), zap.String("email", email))
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, err
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Log.Error("Failed to get user by ID", zap.Error(err), zap.Uint("id", id))
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, err
}
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		logger.Log.Error("Failed to update user", zap.Error(err), zap.Uint("id", user.ID))
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
func (r *userRepository) UpdatePassword(ctx context.Context, id uint, newPassword string) error {
	if id == 0 {
		return fmt.Errorf("invalid user ID")
	}
	if newPassword == "" {
		return fmt.Errorf("password cannot be empty")
	}

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("password", newPassword).
		Error

	if err != nil {
		logger.Log.Error("failed to update password", zap.Error(err), zap.Uint("id", id))
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
