package request

type CreateSeatRequest struct {
	CinemaHallID uint   `json:"cinema_hall_id" binding:"required"`
	SeatNumber   string `json:"seat_number" binding:"required"`
	Row          string `json:"row" binding:"required"`
}
type UpdateSeatRequest struct {
	SeatNumber string `json:"seat_number"`
	Row        string `json:"row"`
	Status     string `json:"status" binding:"omitempty,oneof=available reserved maintenance broken"`
}
