package mapper

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/users/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
)

func ToUserFromRegister(req request.RegisterRequest, hashedPassword string) *model.User {
	return &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     model.RoleUser, // default role
	}
}

func ApplyUpdateProfile(user *model.User, req request.UpdateProfileRequest) {
	user.Name = req.Name
}

func ApplyChangePassword(user *model.User, hashedNewPassword string) {
	user.Password = hashedNewPassword
}

// User AuthResponse (untuk response setelah login / register)
func ToAuthResponse(user *model.User, token string) response.AuthResponse {
	return response.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Token: token,
	}
}
