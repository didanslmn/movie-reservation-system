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
	// Gabung user & admin dalam satu route
	combinedRoutes := protected.Group("/")
	combinedRoutes.Use(middleware.RoleBasedAccess(model.RoleUser, model.RoleAdmin))
	{
		combinedRoutes.GET("/profile", h.Profile)
		combinedRoutes.PUT("/profile", h.UpdateProfile)
		combinedRoutes.PUT("/change-password", h.ChangePassword)
	}

}
