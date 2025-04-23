package handler

import (
	"net/http"
	"strconv"

	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/service"
	"github.com/gin-gonic/gin"
)

type CinemaHallHandler struct {
	service service.CinemaHallService
}

func NewCinemaHallHandler(service service.CinemaHallService) *CinemaHallHandler {
	return &CinemaHallHandler{service: service}
}

func (h *CinemaHallHandler) CreateCinemmaHall(c *gin.Context) {
	var req request.CreateCinemaHallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cinemaHall, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cinemaHall)
}

func (h *CinemaHallHandler) GetAllCinemaHall(c *gin.Context) {
	cinemaHalls, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cinemaHalls)
}

func (h *CinemaHallHandler) GetCinemaHallByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_ID",
			"message": "Invalid cinema hall ID format",
		})
		return
	}
	cinemaHall, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cinemaHall)
}

func (h *CinemaHallHandler) UpdateCinemaHall(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid ID",
			"message": "invalid CinemaHall ID format",
		})
		return
	}
	var req request.UpdateCinemaHallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cinemaHall, err := h.service.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cinemaHall)
}
