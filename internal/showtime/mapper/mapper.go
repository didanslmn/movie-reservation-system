package mapper

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/showtime/model"
)

func ToShowtimeResponse(s *model.Showtime) *response.ShowtimeResponse {
	return &response.ShowtimeResponse{
		ID:           s.ID,
		MovieID:      s.MovieID,
		CinemaHallID: s.CinemaHallID,
		StartTime:    s.StartTime,
		EndTime:      s.EndTime,
	}
}
