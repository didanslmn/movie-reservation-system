package mapper

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/dto/request"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/cinemahall/model"
)

func ToCinemaHallResponse(ch *model.CinemaHall) *response.CinemaHallResponse {
	return &response.CinemaHallResponse{
		ID:       ch.ID,
		Name:     ch.Name,
		Capacity: ch.Capacity,
	}
}
func ToCinemaHallResponses(halls []model.CinemaHall) []response.CinemaHallResponse {
	responses := make([]response.CinemaHallResponse, 0, len(halls))
	for _, hall := range halls {
		responses = append(responses, *ToCinemaHallResponse(&hall))
	}
	return responses
}

func ToCinemaHallModel(req *request.CreateCinemaHallRequest) *model.CinemaHall {
	return &model.CinemaHall{
		Name:     req.Name,
		Capacity: req.Capacity,
	}
}
