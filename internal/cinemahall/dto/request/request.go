package request

type CreateCinemaHallRequest struct {
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity" binding:"required"`
}
type UpdateCinemaHallRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}
