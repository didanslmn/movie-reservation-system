package mapper

import (
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/dto/response"
	"github.com/didanslmn/movie-reservation-system.git/internal/genre/model"
)

func ToGenreResponse(g *model.Genre) *response.Genre {
	return &response.Genre{
		ID:   g.ID,
		Name: g.Name,
	}
}

func ToGenreResponseList(genres []model.Genre) []response.Genre {
	res := make([]response.Genre, len(genres))
	for i, g := range genres {
		res[i] = *ToGenreResponse(&g)
	}
	return res
}
