package model

import (
	"time"

	genreModel "github.com/didanslmn/movie-reservation-system.git/internal/genre/model"
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	Title       string              `gorm:"size:255;not null;index"`
	Description string              `gorm:"type:text"`
	Duration    int                 `gorm:"not null;check:duration >= 30"`
	ReleaseDate time.Time           `gorm:"not null"`
	ImageURL    string              `gorm:"size:512"`
	Rating      float64             `gorm:"type:decimal(3,1);check:rating >= 0 AND rating <= 5"`
	Genres      []*genreModel.Genre `gorm:"many2many:movie_genres;constraint:OnDelete:CASCADE;"`
}
