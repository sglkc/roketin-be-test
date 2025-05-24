/* https://go.dev/doc/tutorial/web-service-gin
 * • API to create and upload movies. Required information related with a movies are at least
 *   title, description, duration, artists, genres
 * • API to update movie
 * • API to list all movies with pagination
 * • API to search movie by title/description/artists/genres
 */
package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/sglkc/roketin-be-test/chal-2/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Movies API
//	@version		1.0
//	@description	This is a sample movies API using Gin framework.
//	@contact.name	sglkc
//	@contact.url	https://github.com/sglkc/roketin-be-test

//	@produce	json
//	@accept		json

// https://gin-gonic.com/en/docs/examples/binding-and-validation/
// https://pkg.go.dev/github.com/go-playground/validator/v10
type Movie struct {
	Id          int      `json:"id" binding:"min=1"`
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Duration    int      `json:"duration" binding:"required,min=1"`
	Artists     []string `json:"artists" binding:"required,min=1"`
	Genres      []string `json:"genres" binding:"required,min=1"`
}

var movies = []Movie{
	{
		Id:          1,
		Title:       "Final Destination: Bloodlines",
		Description: "Plagued by a recurring violent nightmare, a college student returns home to find the one person who can break the cycle and save her family from the horrific fate that inevitably awaits them.",
		Duration:    90,
		Artists:     []string{"Kaitlyn Santa Juana", "Teo Briones", "Rya Kihlstedt"},
		Genres:      []string{"Horror", "Splatter Horror"},
	},
	{
		Id:          2,
		Title:       "Mission: Impossible - The Final Reckoning",
		Description: "Our lives are the sum of our choices. Tom Cruise is Ethan Hunt in Mission: Impossible - The Final Reckoning.",
		Duration:    169,
		Artists:     []string{"Tom Cruise", "Haylett Atwell", "Ving Rhames"},
		Genres:      []string{"Action", "Adeventure", "Thriller"},
	},
}

var idPrimaryKey int = 0

func paginate(c *gin.Context, movies []Movie) []Movie {
	page := c.Query("page")
	limit := c.Query("limit")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}

	start := (pageInt - 1) * limitInt
	end := start + limitInt

	if start > len(movies) {
		return []Movie{}
	}
	if end > len(movies) {
		end = len(movies)
	}

	return movies[start:end]
}

// find movie by id also get latest primary key while at it
func findMovieById(id int) *Movie {
	for _, movie := range movies {
		if movie.Id > idPrimaryKey {
			idPrimaryKey = movie.Id
		}

		if movie.Id == id {
			return &movie
		}
	}
	return nil
}

// https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format

// @Summary		Search movies
// @Description	Search for movies by title, description, artist, or genre
// @Tags			Movies
// @Param			title		query	string	false	"Movie title to search for"
// @Param			description	query	string	false	"Movie description to search for"
// @Param			artist		query	string	false	"Movie artist to search for"
// @Param			genre		query	string	false	"Movie genre to search for"
// @Success		200			{array}	Movie
// @Router			/movies/search [get]
func searchMovie(c *gin.Context) {
	title := strings.ToLower(c.Query("title"))
	description := strings.ToLower(c.Query("description"))
	artist := strings.ToLower(c.Query("artist"))
	genre := strings.ToLower(c.Query("genre"))

	var filteredMovies []Movie

	for _, movie := range movies {
		movieTitle := strings.ToLower(movie.Title)
		movieDescription := strings.ToLower(movie.Description)
		movieArtists := strings.ToLower(strings.Join(movie.Artists, ", "))
		movieGenres := strings.ToLower(strings.Join(movie.Genres, ", "))

		if (strings.Contains(movieTitle, title)) ||
			(strings.Contains(movieDescription, description)) ||
			(strings.Contains(movieArtists, artist)) ||
			(strings.Contains(movieGenres, genre)) {
			filteredMovies = append(filteredMovies, movie)
		}
	}

	c.IndentedJSON(http.StatusOK, paginate(c, filteredMovies))
}

// @Summary		Get all movies
// @Description	Get a list of all movies with pagination
// @Tags			Movies
// @Param			page	query	int	false	"Page number for pagination"	default(1)
// @Param			limit	query	int	false	"Number of movies per page"		default(10)
// @Success		200		{array}	Movie
// @Router			/movies [get]
func getMovies(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, paginate(c, movies))
}

// @Summary		Get movie
// @Description	Get movie by ID
// @Tags			Movies
// @Param			id	path		int	true	"Movie ID"
// @Success		200	{array}		Movie
// @Failure		404	{object}	object{message=string}
// @Router			/movies/{id} [get]
func getMovieById(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid movie ID"})
		return
	}

	movie := findMovieById(idInt)
	if movie != nil {
		c.IndentedJSON(http.StatusOK, *movie)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Movie not found"})
}

// @Summary		Create a new movie
// @Description	Create a new movie
// @Tags			Movies
// @Param			movie	body		Movie	true	"Movie object to create"
// @Success		201		{object}	Movie
// @Failure		400		{object}	object{message=string}
// @Router			/movies [post]
func postMovie(c *gin.Context) {
	var newMovie Movie

	if err := c.ShouldBindJSON(&newMovie); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// assume primary key is the latest ID in the list
	findMovieById(newMovie.Id)
	idPrimaryKey++
	newMovie.Id = idPrimaryKey

	movies = append(movies, newMovie)
	c.IndentedJSON(http.StatusCreated, newMovie)
}

// @Summary		Update a movie
// @Description	Update a movie by ID
// @Tags			Movies
// @Param			id		path		int		true	"Movie ID"
// @Param			movie	body		Movie	true	"Updated movie object"
// @Success		200		{object}	Movie
// @Failure		400		{object}	object{message=string}
// @Failure		404		{object}	object{message=string}
// @Router			/movies/{id} [put]
func updateMovie(c *gin.Context) {
	id := c.Param("id")
	var updatedMovie Movie

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid movie ID"})
		return
	}

	if err := c.ShouldBindJSON(&updatedMovie); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	movie := findMovieById(idInt)
	if movie == nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Movie not found"})
		return
	}

	// check if id is updated, if so, check if it already exists
	if updatedMovie.Id != movie.Id {
		if findMovieById(updatedMovie.Id) != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Movie with updated ID already exists"})
			return
		}
	}

	*movie = updatedMovie
	c.IndentedJSON(http.StatusOK, updatedMovie)
}

// @Summary		Delete a movie
// @Description	Delete a movie by ID
// @Tags			Movies
// @Param			id	path		int	true	"Movie ID"
// @Success		200	{object}	object{message=string}
// @Failure		404	{object}	object{message=string}
// @Router			/movies/{id} [delete]
func deleteMovie(c *gin.Context) {
	id := c.Param("id")

	for i, movie := range movies {
		if id == strconv.Itoa(movie.Id) {
			movies = append(movies[:i], movies[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "movie deleted"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}

func main() {
	for _, movie := range movies {
		if movie.Id > idPrimaryKey {
			idPrimaryKey = movie.Id
		}
	}
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovieById)
	router.GET("/movies/search", searchMovie)
	router.POST("/movies", postMovie)
	router.PUT("/movies/:id", updateMovie)
	router.DELETE("/movies/:id", deleteMovie)

	log.Println("Running at localhost:8080 (docs at http://localhost:8080/swagger/index.html)")
	router.Run("localhost:8080")
}
