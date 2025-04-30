package router

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/middleware"
	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/handler"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/gin-gonic/gin"
)

func ShowtimeRoutes(rg *gin.RouterGroup, h *handler.ShowtimeHandler, jwtSecret string) {
	showtime := rg.Group("/showtimes")

	publicRoutes := showtime.Group("/")
	publicRoutes.Use(middleware.JWTAuthMiddleware(jwtSecret))
	publicRoutes.Use(middleware.RoleBasedAccess(model.RoleUser, model.RoleAdmin))
	{
		publicRoutes.GET("/", h.GetAllShowtimes)
		publicRoutes.GET("/:id", h.GetShowtimeByID)
	}

	// Admin-only
	adminRoutes := showtime.Group("/")
	adminRoutes.Use(middleware.JWTAuthMiddleware(jwtSecret))
	adminRoutes.Use(middleware.RoleBasedAccess(model.RoleAdmin))
	{
		adminRoutes.POST("/", h.CreateShowtime)
		adminRoutes.PUT("/:id", h.UpdateShowtime)
		adminRoutes.DELETE("/:id", h.DeleteShowtime)
	}
}
