package handler

import (
	"net/http"
	"strconv"

	"github.com/didanslmn/movie-reservation-system.git/internal/movie/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/service"
	"github.com/didanslmn/movie-reservation-system.git/utils"
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
	if !utils.BindAndValidate(c, &req) {
		utils.RespondWithError(c, http.StatusBadRequest, "validate failed", nil)
		return
	}

	movie, err := h.service.CreateMovie(c.Request.Context(), req)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Error creating movie", err)
		return
	}

	utils.RespondWithSuccess(c, "Movie created successfully", movie)
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	movie, err := h.service.GetMovie(c.Request.Context(), uint(id))
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Movie not found", err)
		return
	}

	utils.RespondWithSuccess(c, "Movie fetched successfully", movie)
}

func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	movies, err := h.service.GetAllMovies(c.Request.Context())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get movies", err)
		return
	}

	utils.RespondWithSuccess(c, "Movies fetched successfully", movies)
}

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	var req request.UpdateMovie
	if !utils.BindAndValidate(c, &req) {
		utils.RespondWithError(c, http.StatusBadRequest, "validate failed", nil)
		return
	}

	movie, err := h.service.UpdateMovie(c.Request.Context(), uint(id), req)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Error updating movie", err)
		return
	}

	utils.RespondWithSuccess(c, "Movie updated successfully", movie)
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	if err := h.service.DeleteMovie(c.Request.Context(), uint(id)); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Error deleting movie", err)
		return
	}

	utils.RespondWithSuccess(c, "Movie deleted successfully", nil)
}

func (h *MovieHandler) GetMoviesByGenre(c *gin.Context) {
	genreID, err := strconv.ParseUint(c.Param("genre_id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid genre ID", err)
		return
	}

	movies, err := h.service.GetMoviesByGenre(c.Request.Context(), uint(genreID))
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get movies by genre", err)
		return
	}

	utils.RespondWithSuccess(c, "Movies by genre fetched successfully", movies)
}
