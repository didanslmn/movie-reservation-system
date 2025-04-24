package response

import "github.com/didanslmn/movie-reservation-system.git/internal/users/model"

type AuthResponse struct {
	ID    uint       `json:"id"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
	Role  model.Role `json:"role"`
	Token string     `json:"token"`
}
