package handler

import (
	"net/http"

	"github.com/didanslmn/movie-reservation-system.git/internal/middleware"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/service"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if !utils.BindAndValidate(c, &req) {
		utils.ErrorLogger.Printf("Register: Validation failed for email: %s", req.Email)
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	res, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		utils.ErrorLogger.Printf("Register: Failed to register user (email: %s): %v", req.Email, err)
		utils.RespondWithError(c, http.StatusBadRequest, "Failed to register user", err)
		return
	}

	utils.InfoLogger.Printf("Register: User registered successfully (email: %s)", req.Email)
	utils.RespondWithSuccess(c, "User registered successfully", res)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if !utils.BindAndValidate(c, &req) {
		utils.ErrorLogger.Printf("Login: Validation failed for email: %s", req.Email)
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	res, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		utils.ErrorLogger.Printf("Login: Failed login attempt for email: %s: %v", req.Email, err)
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	utils.InfoLogger.Printf("Login: User logged in successfully (email: %s)", req.Email)
	utils.RespondWithSuccess(c, "Login successful", res)
}

func (h *UserHandler) Profile(c *gin.Context) {
	user, ok := middleware.GetUserFromContext(c.Request.Context())
	if !ok || user == nil {
		utils.ErrorLogger.Println("Profile: User not authenticated")
		utils.RespondWithError(c, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	utils.InfoLogger.Printf("Profile: User profile accessed (userID: %d)", user.ID)
	utils.RespondWithSuccess(c, "User profile fetched successfully", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	user, ok := middleware.GetUserFromContext(c.Request.Context())
	if !ok || user == nil {
		utils.ErrorLogger.Println("UpdateProfile: User not authenticated")
		utils.RespondWithError(c, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	var req request.UpdateProfileRequest
	if !utils.BindAndValidate(c, &req) {
		utils.ErrorLogger.Printf("UpdateProfile: Validation failed (userID: %d)", user.ID)
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	res, err := h.userService.UpdateProfile(c.Request.Context(), user.ID, req)
	if err != nil {
		utils.ErrorLogger.Printf("UpdateProfile: Failed to update profile (userID: %d): %v", user.ID, err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	utils.InfoLogger.Printf("UpdateProfile: Profile updated successfully (userID: %d)", user.ID)
	utils.RespondWithSuccess(c, "Profile updated successfully", res)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	user, ok := middleware.GetUserFromContext(c.Request.Context())
	if !ok || user == nil {
		utils.ErrorLogger.Println("ChangePassword: User not authenticated")
		utils.RespondWithError(c, http.StatusUnauthorized, "Authentication required", nil)
		return
	}

	var req request.ChangePasswordRequest
	if !utils.BindAndValidate(c, &req) {
		utils.ErrorLogger.Printf("ChangePassword: Validation failed (userID: %d)", user.ID)
		utils.RespondWithError(c, http.StatusBadRequest, "Validation failed", nil)
		return
	}

	err := h.userService.ChangePassword(c.Request.Context(), user.ID, req)
	if err != nil {
		utils.ErrorLogger.Printf("ChangePassword: Failed to change password (userID: %d): %v", user.ID, err)
		utils.RespondWithError(c, http.StatusBadRequest, "Failed to change password", err)
		return
	}

	utils.InfoLogger.Printf("ChangePassword: Password changed successfully (userID: %d)", user.ID)
	utils.RespondWithSuccess(c, "Password changed successfully", nil)
}
