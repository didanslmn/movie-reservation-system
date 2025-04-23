package router

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/handler"
	"github.com/didanslmn/movie-reservation-system.git/internal/middleware"
	"github.com/didanslmn/movie-reservation-system.git/internal/users/model"
	"github.com/gin-gonic/gin"
)

func GenreRoutes(rg *gin.RouterGroup, h *handler.GenreHandler, jwtSecret string) {
	genre := rg.Group("/genre")

	adminRoutes := genre.Group("/")
	adminRoutes.Use(middleware.JWTAuthMiddleware(jwtSecret))
	adminRoutes.Use(middleware.RoleBasedAccess(model.RoleAdmin))
	{
		adminRoutes.POST("/", h.CreateGenre)
		adminRoutes.PUT("/:id", h.UpdateGenre)
		adminRoutes.DELETE("/:id", h.DeleteGenre)
	}

	publcRoutes := genre.Group("/")
	publcRoutes.Use(middleware.JWTAuthMiddleware(jwtSecret))
	publcRoutes.Use(middleware.RoleBasedAccess(model.RoleAdmin, model.RoleUser))
	{
		publcRoutes.GET("/:id", h.GetGenre)
		publcRoutes.GET("/", h.GetAllGenres)
	}

}
