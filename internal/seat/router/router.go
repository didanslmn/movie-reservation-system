package router

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/middleware"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/handler"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/gin-gonic/gin"
)

func SeatRouts(rg *gin.RouterGroup, h *handler.SeatHandler, jwtSecret string) {
	seat := rg.Group("/seat")
	adminRouts := seat.Group("/")
	adminRouts.Use(middleware.JWTAuthMiddleware(jwtSecret))
	adminRouts.Use(middleware.RoleBasedAccess(model.RoleAdmin))
	{
		adminRouts.POST("/", h.CreateSeat)
		adminRouts.GET("/:id", h.GetSeat)
		adminRouts.GET("/cinemahall/:hall_id", h.GetSeatsByHallID)
		adminRouts.PUT("/:id", h.UpdateSeat)
		adminRouts.DELETE("/:id", h.DeleteSeat)
	}
}
