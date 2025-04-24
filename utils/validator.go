package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func BindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return false
	}

	if err := validate.Struct(req); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Validation failed", err)
		return false
	}

	return true
}
