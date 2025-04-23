package mapper

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/seat/model"
)

func ToSeatResponse(s *model.Seat) *response.SeatResponse {
	return &response.SeatResponse{
		ID:           s.ID,
		SeatNumber:   s.SeatNumber,
		Row:          s.Row,
		Status:       s.Status,
		CinemaHallID: s.CinemaHallID,
	}
}
func ToSeatResponseList(seats []model.Seat) []response.SeatResponse {
	responses := make([]response.SeatResponse, 0, len(seats))
	for i := range seats {
		responses = append(responses, *ToSeatResponse(&seats[i]))
	}
	return responses
}

func ToSeatModel(req *request.CreateSeatRequest) *model.Seat {
	return &model.Seat{
		CinemaHallID: req.CinemaHallID,
		SeatNumber:   req.SeatNumber,
		Row:          req.Row,
		Status:       "available", // default
	}
}
