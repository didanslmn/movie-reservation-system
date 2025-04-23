package handler

import (
	"net/http"
	"strconv"

	"github.com/didanslmn/movie-reservation-system.git/internal/movie/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/service"
	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	service service.MovieService
}

func NewMovieHandler(service service.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req request.CreateMovie
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	movie, err := h.service.CreateMovie(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	movie, err := h.service.GetMovie(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	movies, err := h.service.GetAllMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var req request.UpdateMovie
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	movie, err := h.service.UpdateMovie(c.Request.Context(), uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	if err := h.service.DeleteMovie(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *MovieHandler) GetMoviesByGenre(c *gin.Context) {
	genreID, err := strconv.ParseUint(c.Param("genre_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	movies, err := h.service.GetMoviesByGenre(c.Request.Context(), uint(genreID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movies by genre"})
		return
	}

	c.JSON(http.StatusOK, movies)
}
