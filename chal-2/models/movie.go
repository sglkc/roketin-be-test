package models

// https://gin-gonic.com/en/docs/examples/binding-and-validation/
// https://pkg.go.dev/github.com/go-playground/validator/v10
type Movie struct {
	Id          int      `json:"id"`
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Duration    int      `json:"duration" binding:"required,min=1"`
	Artists     []string `json:"artists" binding:"required,min=1"`
	Genres      []string `json:"genres" binding:"required,min=1"`
}

type Movies []Movie
