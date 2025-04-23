package handler

import (
	"net/http"
	"strconv"

	"github.com/didanslmn/movie-reservation-system.git/internal/seat/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/service"
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reuest", "details": err.Error()})
		return
	}
	seat, err := h.service.CreateSeat(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, seat)
}

func (h *SeatHandler) GetSeat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seat ID"})
		return
	}

	seat, err := h.service.GetSeat(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seat)
}
func (h *SeatHandler) GetSeatsByHallID(c *gin.Context) {
	hallID, err := strconv.ParseUint(c.Param("hall_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hall ID"})
		return
	}

	seats, err := h.service.GetByHallID(c.Request.Context(), uint(hallID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seats)
}
func (h *SeatHandler) UpdateSeat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seat ID"})
		return
	}

	var req request.UpdateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seat, err := h.service.UpdateSeat(c.Request.Context(), uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seat)
}
func (h *SeatHandler) DeleteSeat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seat ID"})
		return
	}

	if err := h.service.DeleteSeat(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
