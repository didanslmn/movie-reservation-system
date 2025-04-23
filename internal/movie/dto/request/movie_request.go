package request

import "time"

type CreateMovie struct {
	Title       string    `json:"title" binding:"required,min=1,max=100"`
	Description string    `json:"description" binding:"required,max=500"`
	Duration    int       `json:"duration" binding:"required,min=30"`
	ReleaseDate time.Time `json:"release_date" binding:"required"`
	ImageURL    string    `json:"image_url" binding:"omitempty,url"`
	Rating      float64   `json:"rating" binding:"numeric,min=0,max=5"`
	GenreIDs    []uint    `json:"genre_ids" binding:"required,min=1"`
}
type UpdateMovie struct {
	Title       string    `json:"title" binding:"omitempty,min=2,max=100"`
	Description string    `json:"description" binding:"omitempty,max=500"`
	Duration    int       `json:"duration" binding:"omitempty,min=30"`
	ReleaseDate time.Time `json:"release_date" binding:"omitempty"`
	ImageURL    string    `json:"image_url" binding:"omitempty,url"`
	Rating      float64   `json:"rating" binding:"omitempty,numeric,min=0,max=5"`
	GenreIDs    []uint    `json:"genre_ids" binding:"omitempty,min=1"`
}
