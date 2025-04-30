package router

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/middleware"
	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/handler"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/gin-gonic/gin"
)

func ReservationRoutes(rg *gin.RouterGroup, h *handler.ReservationHandler, jwtSecret string) {
	reservations := rg.Group("/reservations")

	// User & Admin dapat membuat, melihat semua, dan melihat detail
	publicRoutes := reservations.Group("/")
	publicRoutes.Use(middleware.JWTAuthMiddleware(jwtSecret))
	publicRoutes.Use(middleware.RoleBasedAccess(model.RoleUser, model.RoleAdmin))
	{
		publicRoutes.POST("/", h.CreateReservation)
		publicRoutes.GET("/", h.GetAllReservations)
		publicRoutes.GET("/:id", h.GetReservationByID)
	}

	// Hanya Admin yang dapat menghapus reservation
	adminRoutes := reservations.Group("/")
	adminRoutes.Use(middleware.JWTAuthMiddleware(jwtSecret))
	adminRoutes.Use(middleware.RoleBasedAccess(model.RoleAdmin))
	{
		adminRoutes.DELETE("/:id", h.DeleteReservation)
	}
}
