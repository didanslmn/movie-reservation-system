package response

import "time"

type ShowtimeResponse struct {
	ID           uint      `json:"id"`
	MovieID      uint      `json:"movie_id"`
	CinemaHallID uint      `json:"cinema_hall_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
}
