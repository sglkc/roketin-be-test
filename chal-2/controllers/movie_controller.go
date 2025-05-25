package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sglkc/roketin-be-test/chal-2/database"
	"github.com/sglkc/roketin-be-test/chal-2/models"
	"github.com/sglkc/roketin-be-test/chal-2/utils"
)

// find movie by id also get latest primary key while at it
func findMovieById(id int) *models.Movie {
	for _, movie := range database.Movies {
		if movie.Id > database.MovieId {
			database.MovieId = movie.Id
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
func SearchMovie(c *gin.Context) {
	title := strings.ToLower(c.Query("title"))
	description := strings.ToLower(c.Query("description"))
	artist := strings.ToLower(c.Query("artist"))
	genre := strings.ToLower(c.Query("genre"))

	var filteredMovies models.Movies

	for _, movie := range database.Movies {
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

	c.IndentedJSON(http.StatusOK, utils.Paginate(c, filteredMovies))
}

// @Summary		Get all movies
// @Description	Get a list of all movies with pagination
// @Tags			Movies
// @Param			page	query	int	false	"Page number for pagination"	default(1)
// @Param			limit	query	int	false	"Number of movies per page"		default(10)
// @Success		200		{array}	Movie
// @Router			/movies [get]
func GetMovies(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, utils.Paginate(c, database.Movies))
}

// @Summary		Get movie
// @Description	Get movie by ID
// @Tags			Movies
// @Param			id	path		int	true	"Movie ID"
// @Success		200	{array}		Movie
// @Failure		404	{object}	object{message=string}
// @Router			/movies/{id} [get]
func GetMovieById(c *gin.Context) {
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
func PostMovie(c *gin.Context) {
	var newMovie models.Movie

	if err := c.ShouldBindJSON(&newMovie); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// assume primary key is the latest ID in the list
	findMovieById(newMovie.Id)
	database.MovieId++
	newMovie.Id = database.MovieId

	database.Movies = append(database.Movies, newMovie)
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
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	var updatedMovie models.Movie

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
func DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	for i, movie := range database.Movies {
		if id == strconv.Itoa(movie.Id) {
			database.Movies = append(database.Movies[:i], database.Movies[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "movie deleted"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}
