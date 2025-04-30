package response

import "time"

type ReservationResponse struct {
	ID       uint             `json:"id"`
	User     UserResponse     `json:"user"`
	Showtime ShowtimeResponse `json:"showtime"`
	Seats    []SeatResponse   `json:"seat"`
	Status   string           `json:"status"`
}

type UserResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ShowtimeResponse struct {
	ID        uint      `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type SeatResponse struct {
	ID         uint   `json:"id"`
	SeatNumber string `json:"seat_number"`
	Row        string `json:"row"`
}
