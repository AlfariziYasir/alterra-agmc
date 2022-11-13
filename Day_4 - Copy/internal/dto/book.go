package dto

type (
	BookRequest struct {
		ID     uint   `param:"id"`
		Title  string `json:"title" validate:"required"`
		Isbn   string `json:"isbn" validate:"required"`
		Writer string `json:"writer" validate:"required"`
	}

	BookResponse struct {
		ID     uint   `json:"id"`
		Title  string `json:"title"`
		Isbn   string `json:"isbn"`
		Writer string `json:"writer"`
	}
)
