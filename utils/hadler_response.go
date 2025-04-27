package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func RespondWithError(c *gin.Context, code int, msg string, err error) {
	ErrorLogger.Printf("%s: %v", msg, err)
	c.JSON(code, ErrorResponse{
		Message: msg,
		Error:   err.Error(),
	})
}

func RespondWithSuccess(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: msg,
		Data:    data,
	})
}
