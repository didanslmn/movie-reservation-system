package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/didanslmn/movie-reservation-system.git/utils"
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
		utils.ErrorLogger.Printf("Create user failed (email: %s): %v", user.Email, err)
		return fmt.Errorf("failed to create user: %v", err)
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
		utils.ErrorLogger.Printf("Get user by email failed (email: %s): %v", email, err)
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		utils.ErrorLogger.Printf("Get user by ID failed (id: %d): %v", id, err)
		return nil, fmt.Errorf("failed to get user by ID: %v", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		utils.ErrorLogger.Printf("Update user failed (id: %d): %v", user.ID, err)
		return fmt.Errorf("failed to update user: %v", err)
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
		utils.ErrorLogger.Printf("Update password failed (id: %d): %v", id, err)
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}
