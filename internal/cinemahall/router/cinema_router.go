package router

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/handler"
	"github.com/didanslmn/movie-reservation-system.git/internal/middleware"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/gin-gonic/gin"
)

func CinemaHallRouts(rg *gin.RouterGroup, h *handler.CinemaHallHandler, jwtSecret string) {
	cinemas := rg.Group("/cinemahalls")
	adminRouts := cinemas.Group("/")
	adminRouts.Use(middleware.JWTAuthMiddleware(jwtSecret))
	adminRouts.Use(middleware.RoleBasedAccess(model.RoleAdmin))
	{
		adminRouts.POST("/", h.CreateCinemaHall)
		adminRouts.GET("/", h.GetAllCinemaHall)
		adminRouts.GET("/:id", h.GetCinemaHallByID)
		adminRouts.PUT("/:id", h.UpdateCinemaHall)
	}
}
