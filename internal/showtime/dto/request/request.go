package request

import "time"

type CreateShowtimeRequest struct {
	MovieID      uint      `json:"movie_id" binding:"required"`
	CinemaHallID uint      `json:"cinema_hall_id" binding:"required"`
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
}

type UpdateShowtimeRequest struct {
	MovieID      uint      `json:"movie_id" binding:"required"`
	CinemaHallID uint      `json:"cinema_hall_id" binding:"required"`
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
}
