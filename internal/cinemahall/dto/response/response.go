package response

type CinemaHallResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}
