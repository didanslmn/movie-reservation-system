package handler

import (
	"net/http"

	"github.com/didanslmn/movie-reservation-system.git/internal/middleware"
	"github.com/didanslmn/movie-reservation-system.git/internal/pkg/logger"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/service"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}
func (h *UserHandler) bindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		logger.Log.Debug("Failed to bind request", zap.Error(err))
		h.sendValidationError(c, err)
		return false
	}

	if err := h.validate.Struct(req); err != nil {
		logger.Log.Debug("Validation failed", zap.Error(err))
		h.sendValidationError(c, err)
		return false
	}
	return true
}

func (h *UserHandler) sendValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[e.Field()] = h.getValidationMessage(e)
		}
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "validation_error",
			Message: "Validation failed",
			Errors:  errors,
		})
		return
	}

	c.JSON(http.StatusBadRequest, response.ErrorResponse{
		Code:    "invalid_request",
		Message: err.Error(),
	})
}

func (h *UserHandler) getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Should be at least " + e.Param() + " characters"
	}
	return "Invalid value"
}

func (h *UserHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if !h.bindAndValidate(c, &req) {
		return
	}

	res, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		logger.Log.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "registration_failed",
			Message: err.Error()})
		return
	}
	logger.Log.Info("user register succesfuly",
		zap.String("email", req.Email),
		zap.String("name", req.Name))

	c.JSON(http.StatusCreated, response.SuccessResponse{
		Data: res,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if !h.bindAndValidate(c, &req) {
		return
	}

	res, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		logger.Log.Warn("Login failed",
			zap.String("email", req.Email),
			zap.Error(err))
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Code:    "authentication_failed",
			Message: "Invalid credentials",
		})
		return
	}
	logger.Log.Info("User logged in",
		zap.String("email", req.Email))

	c.JSON(http.StatusOK, response.SuccessResponse{
		Data: res,
	})
}
func (h *UserHandler) Profile(c *gin.Context) {
	user, ok := middleware.GetUserFromContext(c.Request.Context())
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Code:    "unauthorized",
			Message: "Authentication required",
		})
		return
	}

	c.JSON(200, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *UserHandler) AdminDashboard(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Welcome, Admin!"})
}
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	user, ok := middleware.GetUserFromContext(c.Request.Context())
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Code:    "unauthorized",
			Message: "Authentication required",
		})
		return
	}
	var req request.UpdateProfileRequest
	if !h.bindAndValidate(c, &req) {
		return
	}

	// Panggil service untuk update profile
	res, err := h.userService.UpdateProfile(c.Request.Context(), user.ID, req)
	if err != nil {
		logger.Log.Error("Failed to update profile",
			zap.Uint("userID", user.ID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Code:    "update_failed",
			Message: err.Error(),
		})
		return
	}
	logger.Log.Info("User updated profile",
		zap.Uint("userID", user.ID))

	c.JSON(http.StatusOK, response.SuccessResponse{
		Data: res,
	})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	user, ok := middleware.GetUserFromContext(c.Request.Context())
	if !ok || user == nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Code:    "unauthorized",
			Message: "Authentication required",
		})
		return
	}

	var req request.ChangePasswordRequest
	if !h.bindAndValidate(c, &req) {
		return
	}

	err := h.userService.ChangePassword(c.Request.Context(), user.ID, req)
	if err != nil {
		logger.Log.Error("Failed to change password",
			zap.Uint("userID", user.ID),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Code:    "password_change_failed",
			Message: err.Error(),
		})
		return
	}
	logger.Log.Info("User changed password",
		zap.Uint("userID", user.ID))

	c.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Password updated successfully",
	})
}
