package mapper

import (
	genreModel "github.com/didanslmn/movie-reservation-system.git/internal/genre/model"
	"github.com/didanslmn/movie-reservation-system.git/internal/movie/dto/response"
	movieModel "github.com/didanslmn/movie-reservation-system.git/internal/movie/model"
)

func ToMovieResponse(m *movieModel.Movie) *response.Movie {
	return &response.Movie{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Duration:    m.Duration,
		ReleaseDate: m.ReleaseDate,
		ImageURL:    m.ImageURL,
		Rating:      m.Rating,
		Genres:      convertGenres(m.Genres),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func convertGenres(genres []*genreModel.Genre) []response.Genre {
	result := make([]response.Genre, len(genres))
	for i, g := range genres {
		result[i] = response.Genre{
			ID:   g.ID,
			Name: g.Name,
		}
	}
	return result
}
