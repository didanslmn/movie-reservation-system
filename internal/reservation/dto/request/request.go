package request

type CreateReservationRequest struct {
	ShowtimeID uint   `json:"showtime_id" binding:"required"`
	SeatIDs    []uint `json:"seat_id" binding:"required,min=1"`
}
