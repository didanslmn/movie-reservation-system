package response

import "time"

type Movie struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Duration    int       `json:"duration"`
	ReleaseDate time.Time `json:"release_date"`
	ImageURL    string    `json:"image_url"`
	Rating      float64   `json:"rating"`
	Genres      []Genre   `json:"genres"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type Genre struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
