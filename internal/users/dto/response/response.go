package response

import "github.com/didanslmn/movie-reservation-system.git/internal/users/model"

type AuthResponse struct {
	ID    uint       `json:"id"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
	Role  model.Role `json:"role"`
	Token string     `json:"token"`
}
type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}
