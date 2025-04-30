package handler

import (
	"net/http"
	"strconv"

	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/service"
	"github.com/didanslmn/movie-reservation-system.git/utils"
	"github.com/gin-gonic/gin"
)

type ReservationHandler struct {
	service service.ReservationService
}

func NewReservationHandler(service service.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var req request.CreateReservationRequest
	if ok := utils.BindAndValidate(c, &req); !ok {
		return
	}

	// Ambil user ID dari context (hasil autentikasi JWT)
	userIDAny, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "Invalid user ID type", nil)
		return
	}

	res, err := h.service.CreateReservation(c.Request.Context(), userID, &req)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create reservation", err)
		return
	}
	utils.RespondWithSuccess(c, "Reservation created successfully", res)
}

func (h *ReservationHandler) GetReservationByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid reservation ID", err)
		return
	}

	res, err := h.service.GetReservationByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Reservation not found", err)
		return
	}

	utils.RespondWithSuccess(c, "Reservation fetched successfully", res)
}

func (h *ReservationHandler) GetAllReservations(c *gin.Context) {
	res, err := h.service.GetAllReservations(c.Request.Context())
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch reservations", err)
		return
	}
	utils.RespondWithSuccess(c, "Reservations fetched successfully", res)
}

func (h *ReservationHandler) DeleteReservation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid reservation ID", err)
		return
	}

	err = h.service.DeleteReservation(c.Request.Context(), uint(id))
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete reservation", err)
		return
	}
	utils.RespondWithSuccess(c, "Reservation deleted successfully", nil)
}
