package service

import (
	"context"
	"fmt"
	"time"

	"github.com/didanslmn/movie-reservation-system.git/internal/users/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/mapper"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/repository"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"github.com/golang-jwt/jwt/v5"
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
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string) UserService {
	return &userService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *userService) Register(ctx context.Context, req request.RegisterRequest) (*response.AuthResponse, error) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		utils.InfoLogger.Printf("register: email already exists (%s)", req.Email)
		return nil, fmt.Errorf("email already registered: %s", req.Email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorLogger.Printf("register: failed to hash password: %v", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := mapper.ToUserFromRegister(req, string(hashedPassword))

	if err := s.userRepo.Create(ctx, user); err != nil {
		utils.ErrorLogger.Printf("register: failed to create user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	utils.InfoLogger.Printf("user registered successfully (email: %s)", user.Email)
	return s.toAuthResponse(user)
}

func (s *userService) Login(ctx context.Context, req request.LoginRequest) (*response.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		utils.InfoLogger.Printf("login: invalid credentials (email: %s)", req.Email)
		return nil, fmt.Errorf("invalid credentials: %s", req.Email)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.InfoLogger.Printf("login: wrong password (email: %s)", req.Email)
		return nil, fmt.Errorf("invalid credentials: %s", req.Email)
	}

	utils.InfoLogger.Printf("user logged in (email: %s)", user.Email)
	return s.toAuthResponse(user)
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint, req request.UpdateProfileRequest) (*response.AuthResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		utils.ErrorLogger.Printf("update profile: failed to get user (id: %d): %v", userID, err)
		return nil, fmt.Errorf("failed to get user with ID %d: %w", userID, err)
	}

	mapper.ApplyUpdateProfile(user, req)

	if err := s.userRepo.Update(ctx, user); err != nil {
		utils.ErrorLogger.Printf("update profile: failed to update user (id: %d): %v", user.ID, err)
		return nil, fmt.Errorf("failed to update user with ID %d: %w", user.ID, err)
	}

	utils.InfoLogger.Printf("profile updated (id: %d)", user.ID)
	return s.toAuthResponse(user)
}

func (s *userService) ChangePassword(ctx context.Context, userID uint, req request.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		utils.ErrorLogger.Printf("change password: failed to get user (id: %d): %v", userID, err)
		return fmt.Errorf("failed to get user with ID %d: %w", userID, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		utils.InfoLogger.Printf("change password: incorrect old password (id: %d)", userID)
		return fmt.Errorf("incorrect old password for user ID %d", userID)
	}

	newHashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorLogger.Printf("change password: failed to hash new password: %v", err)
		return fmt.Errorf("failed to hash new password: %w", err)
	}
	mapper.ApplyChangePassword(user, string(newHashed))
	if err := s.userRepo.UpdatePassword(ctx, user.ID, string(newHashed)); err != nil {
		utils.ErrorLogger.Printf("change password: failed to update password: %v", err)
		return fmt.Errorf("failed to update password for user ID %d: %w", user.ID, err)
	}

	utils.InfoLogger.Printf("password changed (id: %d)", user.ID)
	return nil
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

func (s *userService) toAuthResponse(user *model.User) (*response.AuthResponse, error) {
	token, err := s.generateToken(user)
	if err != nil {
		utils.ErrorLogger.Printf("failed to generate JWT: %v", err)
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	authResp := mapper.ToAuthResponse(user, token)
	return &authResp, nil
}
