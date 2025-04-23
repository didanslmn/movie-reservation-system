package response

type SeatResponse struct {
	ID           uint   `json:"id"`
	SeatNumber   string `json:"seat_number"`
	Row          string `json:"row"`
	Status       string `json:"status"`
	CinemaHallID uint   `json:"cinema_hall_id"`
}
