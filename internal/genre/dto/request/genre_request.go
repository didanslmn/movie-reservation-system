package request

type CreateGenre struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}

type UpdateGenre struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}
