package handler

import (
	"net/http"
	"strconv"

	"github.com/didanslmn/movie-reservation-system.git/internal/genre/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/service"
	"github.com/gin-gonic/gin"
)

type GenreHandler struct {
	service service.GenreService
}

func NewGenreHandler(service service.GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

func (h *GenreHandler) CreateGenre(c *gin.Context) {
	var req request.CreateGenre
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
	}

	res, err := h.service.CreateGenre(c.Request.Context(), req)
	if err != nil {
		handleGenreError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *GenreHandler) GetGenre(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_ID",
			"message": "Invalid genre ID format",
		})
		return
	}

	res, err := h.service.GetGenre(c.Request.Context(), uint(id))
	if err != nil {
		handleGenreError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *GenreHandler) GetAllGenres(c *gin.Context) {
	genres, err := h.service.GetAllGenres(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "GET_GENRES_FAILED",
			"message": "Failed to get genres",
		})
		return
	}

	c.JSON(http.StatusOK, genres)
}

func (h *GenreHandler) UpdateGenre(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_ID",
			"message": "Invalid genre ID format",
		})
		return
	}

	var req request.UpdateGenre
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "VALIDATION_ERROR",
			"message": err.Error(),
		})
		return
	}

	res, err := h.service.UpdateGenre(c.Request.Context(), uint(id), req)
	if err != nil {
		handleGenreError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *GenreHandler) DeleteGenre(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_ID",
			"message": "Invalid genre ID format",
		})
		return
	}

	if err := h.service.DeleteGenre(c.Request.Context(), uint(id)); err != nil {
		handleGenreError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func handleGenreError(c *gin.Context, err error) {
	switch err.Error() {
	case "genre not found":
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "GENRE_NOT_FOUND",
			"message": err.Error(),
		})
	case "genre already exists", "new genre name already exists":
		c.JSON(http.StatusConflict, gin.H{
			"error":   "GENRE_CONFLICT",
			"message": err.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_ERROR",
			"message": "Something went wrong",
		})
	}
}
