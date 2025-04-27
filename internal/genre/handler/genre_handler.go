package handler

import (
	"net/http"
	"strconv"

	"github.com/didanslmn/movie-reservation-system.git/internal/genre/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/service"
	"github.com/didanslmn/movie-reservation-system.git/utils"
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
	if ok := utils.BindAndValidate(c, &req); !ok {
		utils.RespondWithError(c, http.StatusBadRequest, "validate failed", nil)
		return
	}

	res, err := h.service.CreateGenre(c.Request.Context(), req)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "failed to create genre", err)
		return
	}

	utils.RespondWithSuccess(c, "Genre created successfully", res)
}

func (h *GenreHandler) GetGenre(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid genre ID format", err)
		return
	}

	res, err := h.service.GetGenre(c.Request.Context(), uint(id))
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "genre not found", err)
		return
	}

	utils.RespondWithSuccess(c, "Genre retrieved successfully", res)
}

func (h *GenreHandler) GetAllGenres(c *gin.Context) {
	res, err := h.service.GetAllGenres(c.Request.Context())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get genres", err)
		return
	}

	utils.RespondWithSuccess(c, "Genres retrieved successfully", res)
}

func (h *GenreHandler) UpdateGenre(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid genre ID format", err)
		return
	}

	var req request.UpdateGenre
	if ok := utils.BindAndValidate(c, &req); !ok {
		utils.RespondWithError(c, http.StatusBadRequest, "validate failed", nil)
		return
	}

	res, err := h.service.UpdateGenre(c.Request.Context(), uint(id), req)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "failed to update genre", err)
		return
	}

	utils.RespondWithSuccess(c, "Genre updated successfully", res)
}

func (h *GenreHandler) DeleteGenre(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid genre ID format", err)
		return
	}

	if err := h.service.DeleteGenre(c.Request.Context(), uint(id)); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "failed to delete genre", err)
		return
	}

	c.Status(http.StatusNoContent)
}
