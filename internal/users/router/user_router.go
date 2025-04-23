package router

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/middleware"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/handler"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, h *handler.UserHandler) {
	// Group untuk auth
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
	}
}
func UserRoutes(rg *gin.RouterGroup, h *handler.UserHandler, jwtSecret string) {
	protected := rg.Group("/user")
	protected.Use(middleware.JWTAuthMiddleware(jwtSecret))

	// Routes untuk user biasa
	userRoutes := protected.Group("/")
	userRoutes.Use(middleware.RoleBasedAccess(model.RoleUser))
	{
		userRoutes.GET("/profile", h.Profile)
		userRoutes.PUT("/profile", h.UpdateProfile)
		userRoutes.PUT("/change-password", h.ChangePassword)
	}

	// Routes untuk admin
	adminRoutes := protected.Group("/admin")
	adminRoutes.Use(middleware.RoleBasedAccess(model.RoleAdmin))
	{
		adminRoutes.GET("/", h.AdminDashboard)
		adminRoutes.GET("/profile", h.Profile)
		adminRoutes.PUT("/profile", h.UpdateProfile)
		adminRoutes.PUT("/change-password", h.ChangePassword)
	}
}
