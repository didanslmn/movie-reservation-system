package handler

import (
	"net/http"
	"strconv"

	"github.com/didanslmn/movie-reservation-system.git/internal/seat/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/service"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"github.com/gin-gonic/gin"
)

type SeatHandler struct {
	service service.SeatService
}

func NewSeatHandler(service service.SeatService) *SeatHandler {
	return &SeatHandler{service: service}
}

func (h *SeatHandler) CreateSeat(c *gin.Context) {
	var req request.CreateSeatRequest
	if ok := utils.BindAndValidate(c, &req); !ok {
		utils.RespondWithError(c, http.StatusBadRequest, "validate failed", nil)
		return
	}
	seat, err := h.service.CreateSeat(c.Request.Context(), req)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create seat", err)
		return
	}

	utils.RespondWithSuccess(c, "Seat created successfully", seat)
}

func (h *SeatHandler) GetSeat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid seat ID", err)
		return
	}

	seat, err := h.service.GetSeat(c.Request.Context(), uint(id))
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Seat not found", err)
		return
	}

	utils.RespondWithSuccess(c, "Seat retrieved successfully", seat)
}

func (h *SeatHandler) GetSeatsByHallID(c *gin.Context) {
	hallID, err := strconv.ParseUint(c.Param("hall_id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid hall ID", err)
		return
	}

	seats, err := h.service.GetByHallID(c.Request.Context(), uint(hallID))
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve seats", err)
		return
	}

	utils.RespondWithSuccess(c, "Seats retrieved successfully", seats)
}

func (h *SeatHandler) UpdateSeat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid seat ID", err)
		return
	}

	var req request.UpdateSeatRequest
	if ok := utils.BindAndValidate(c, &req); !ok {
		utils.RespondWithError(c, http.StatusBadRequest, "validate failed", nil)
		return
	}

	seat, err := h.service.UpdateSeat(c.Request.Context(), uint(id), req)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update seat", err)
		return
	}

	utils.RespondWithSuccess(c, "Seat updated successfully", seat)
}

func (h *SeatHandler) DeleteSeat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid seat ID", err)
		return
	}

	if err := h.service.DeleteSeat(c.Request.Context(), uint(id)); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete seat", err)
		return
	}

	c.Status(http.StatusNoContent)
}
