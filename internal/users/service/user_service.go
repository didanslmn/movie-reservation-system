package service

import (
	"context"
	"errors"
	"time"

	"github.com/didanslmn/movie-reservation-system.git/internal/users/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/repository"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req request.RegisterRequest) (*response.AuthResponse, error)
	Login(ctx context.Context, req request.LoginRequest) (*response.AuthResponse, error)
	ChangePassword(ctx context.Context, userID uint, req request.ChangePasswordRequest) error
	UpdateProfile(ctx context.Context, userID uint, req request.UpdateProfileRequest) (*response.AuthResponse, error)
}

type userService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	logger    *zap.Logger
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string, logger *zap.Logger) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *userService) Register(ctx context.Context, req request.RegisterRequest) (*response.AuthResponse, error) {
	// cek user apakah sudah ada
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		s.logger.Warn("register: email already exists", zap.String("email", req.Email))
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("register: failed to hash password", zap.Error(err))
		return nil, err
	}

	// Create user
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     model.RoleUser, // Default role
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("register: failed to create user", zap.Error(err))
		return nil, err
	}
	s.logger.Info("user registered successfully", zap.String("email", user.Email))
	return s.toAuthResponse(user)

}

func (s *userService) Login(ctx context.Context, req request.LoginRequest) (*response.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		s.logger.Warn("login: invalid credentials", zap.String("email", req.Email))
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Warn("login: wrong password", zap.String("email", req.Email))
		return nil, errors.New("invalid credentials")
	}

	s.logger.Info("user logged in", zap.String("email", user.Email))
	return s.toAuthResponse(user)
}

func (s *userService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
func (s *userService) UpdateProfile(ctx context.Context, userID uint, req request.UpdateProfileRequest) (*response.AuthResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("update profile: failed to get user", zap.Uint("userID", userID), zap.Error(err))
		return nil, err
	}

	user.Name = req.Name

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("update profile: failed to update user", zap.Uint("userID", user.ID), zap.Error(err))
		return nil, err
	}

	s.logger.Info("profile updated", zap.Uint("userID", user.ID))
	return s.toAuthResponse(user)
}

func (s *userService) ChangePassword(ctx context.Context, userID uint, req request.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("change password: failed to get user", zap.Uint("userID", userID), zap.Error(err))
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		s.logger.Warn("change password: incorrect old password", zap.Uint("userID", userID))
		return errors.New("incorrect old password")
	}

	newHashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("change password: failed to hash new password", zap.Error(err))
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, string(newHashed)); err != nil {
		s.logger.Error("change password: failed to update password", zap.Error(err))
		return err
	}

	s.logger.Info("password changed", zap.Uint("userID", user.ID))
	return nil
}

func (s *userService) toAuthResponse(user *model.User) (*response.AuthResponse, error) {
	token, err := s.generateToken(user)
	if err != nil {
		s.logger.Error("failed to generate JWT", zap.Error(err))
		return nil, err
	}

	return &response.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Token: token,
	}, nil
}
