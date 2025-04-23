package mapper

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/users/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"go.uber.org/zap"
)

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
