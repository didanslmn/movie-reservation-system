package mapper

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/reservation/model"
)

func ToReservationResponse(r *model.Reservation) *response.ReservationResponse {
	return &response.ReservationResponse{
		ID: r.ID,
		User: response.UserResponse{
			ID:   r.User.ID,
			Name: r.User.Name,
		},
		Showtime: response.ShowtimeResponse{
			ID:        r.Showtime.ID,
			StartTime: r.Showtime.StartTime,
			EndTime:   r.Showtime.EndTime,
		},
		Status: r.Status,
		Seats: func() []response.SeatResponse {
			seats := make([]response.SeatResponse, 0, len(r.Seats))
			for _, seat := range r.Seats {
				seats = append(seats, response.SeatResponse{
					ID:         seat.ID,
					SeatNumber: seat.SeatNumber,
					Row:        seat.Row,
				})
			}
			return seats
		}(),
	}
}

func ToReservationResponseList(reservations []model.Reservation) []response.ReservationResponse {
	responses := make([]response.ReservationResponse, 0, len(reservations))
	for i := range reservations {
		responses = append(responses, *ToReservationResponse(&reservations[i]))
	}
	return responses
}
