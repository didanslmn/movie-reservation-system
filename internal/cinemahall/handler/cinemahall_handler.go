package handler

import (
	"net/http"
	"strconv"

	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/service"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"github.com/gin-gonic/gin"
)

type CinemaHallHandler struct {
	service service.CinemaHallService
}

func NewCinemaHallHandler(service service.CinemaHallService) *CinemaHallHandler {
	return &CinemaHallHandler{service: service}
}

func (h *CinemaHallHandler) CreateCinemaHall(c *gin.Context) {
	var req request.CreateCinemaHallRequest
	if ok := utils.BindAndValidate(c, &req); !ok {
		utils.RespondWithError(c, http.StatusBadRequest, "validate failed", nil)
		return
	}

	cinemaHall, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Failed to create cinema hall", err)
		return
	}

	utils.RespondWithSuccess(c, "Cinema hall created successfully", cinemaHall)
}

func (h *CinemaHallHandler) GetAllCinemaHall(c *gin.Context) {
	cinemaHalls, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve cinema halls", err)
		return
	}

	utils.RespondWithSuccess(c, "Cinema halls retrieved successfully", cinemaHalls)
}

func (h *CinemaHallHandler) GetCinemaHallByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid cinema hall ID format", err)
		return
	}

	cinemaHall, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Failed to retrieve cinema hall", err)
		return
	}

	utils.RespondWithSuccess(c, "Cinema hall retrieved successfully", cinemaHall)
}

func (h *CinemaHallHandler) UpdateCinemaHall(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid cinema hall ID format", err)
		return
	}

	var req request.UpdateCinemaHallRequest
	if ok := utils.BindAndValidate(c, &req); !ok {
		utils.RespondWithError(c, http.StatusBadRequest, "validate failed", nil)
		return
	}

	cinemaHall, err := h.service.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update cinema hall", err)
		return
	}

	utils.RespondWithSuccess(c, "Cinema hall updated successfully", cinemaHall)
}

func (h *CinemaHallHandler) DeleteCinemaHall(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid cinema hall ID format", err)
		return
	}

	err = h.service.Delete(c.Request.Context(), uint(id))
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete cinema hall", err)
		return
	}

	utils.RespondWithSuccess(c, "Cinema hall deleted successfully", nil)
}
