package handler

import (
	"net/http"
	"strconv"

	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/service"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"github.com/gin-gonic/gin"
)

type ShowtimeHandler struct {
	showtimeService service.ShowtimeService
}

func NewShowtimeHandler(showtimeService service.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{showtimeService: showtimeService}
}

func (h *ShowtimeHandler) CreateShowtime(c *gin.Context) {
	var req request.CreateShowtimeRequest
	if !utils.BindAndValidate(c, &req) {
		return
	}

	showtime, err := h.showtimeService.CreateShowtime(c.Request.Context(), req)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create showtime", err)
		return
	}

	utils.RespondWithSuccess(c, "Showtime created successfully", showtime)
}

func (h *ShowtimeHandler) GetShowtimeByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid mshowtimr ID", err)
		return
	}

	showtime, err := h.showtimeService.GetShowtimeByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Showtime not found", err)
		return
	}

	utils.RespondWithSuccess(c, "Showtime fetched successfully", showtime)
}

func (h *ShowtimeHandler) GetAllShowtimes(c *gin.Context) {
	showtimes, err := h.showtimeService.GetAllShowtimes(c.Request.Context())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch showtimes", err)
		return
	}

	utils.RespondWithSuccess(c, "Showtimes fetched successfully", showtimes)
}

func (h *ShowtimeHandler) UpdateShowtime(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid shwotime ID", err)
		return
	}

	var req request.UpdateShowtimeRequest
	if !utils.BindAndValidate(c, &req) {
		return
	}

	showtime, err := h.showtimeService.UpdateShowtime(c.Request.Context(), uint(id), req)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update showtime", err)
		return
	}

	utils.RespondWithSuccess(c, "Showtime updated successfully", showtime)
}

func (h *ShowtimeHandler) DeleteShowtime(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid showtime ID", err)
		return
	}
	err = h.showtimeService.DeleteShowtime(c.Request.Context(), uint(id))
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete showtime", err)
		return
	}

	utils.RespondWithSuccess(c, "Showtime deleted successfully", nil)
}
