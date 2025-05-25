package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sglkc/roketin-be-test/chal-2/database"
	"github.com/sglkc/roketin-be-test/chal-2/dto"
	"github.com/sglkc/roketin-be-test/chal-2/models"
	"github.com/sglkc/roketin-be-test/chal-2/utils"
)

// https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format

// @Summary		Search movies
// @Description	Search for movies by title, description, artist, or genre
// @Tags			Movies
// @Param			title		query	string	false	"Movie title to search for"
// @Param			description	query	string	false	"Movie description to search for"
// @Param			artist		query	string	false	"Movie artist to search for"
// @Param			genre		query	string	false	"Movie genre to search for"
// @Param			page		query	int		false	"Page number for pagination"	default(1)
// @Param			limit		query	int		false	"Number of movies per page"		default(10)
// @Success		200			{array}	dto.PaginatedResponse[models.Movie]
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

		if (title != "" && strings.Contains(movieTitle, title)) ||
			(description != "" && strings.Contains(movieDescription, description)) ||
			(artist != "" && strings.Contains(movieArtists, artist)) ||
			(genre != "" && strings.Contains(movieGenres, genre)) {
			filteredMovies = append(filteredMovies, movie)
		}
	}

	data, page, limit := utils.Paginate(c, filteredMovies)

	c.IndentedJSON(http.StatusOK, dto.PaginatedResponse[models.Movie]{
		BaseResponse: dto.BaseResponse{
			Message: "Movies found",
			Success: true,
		},
		Data:  data,
		Page:  page,
		Limit: limit,
		Count: len(filteredMovies),
	})
}

// @Summary		Get all movies
// @Description	Get a list of all movies with pagination
// @Tags			Movies
// @Param			page	query	int	false	"Page number for pagination"	default(1)
// @Param			limit	query	int	false	"Number of movies per page"		default(10)
// @Success		200		{array}	dto.PaginatedResponse[models.Movie]
// @Router			/movies [get]
func GetMovies(c *gin.Context) {
	data, page, limit := utils.Paginate(c, database.Movies)

	c.IndentedJSON(http.StatusOK, dto.PaginatedResponse[models.Movie]{
		BaseResponse: dto.BaseResponse{
			Message: "Movies found",
			Success: true,
		},
		Data:  data,
		Page:  page,
		Limit: limit,
		Count: len(database.Movies),
	})
}

// @Summary		Get movie
// @Description	Get movie by ID
// @Tags			Movies
// @Param			id	path		int	true	"Movie ID"
// @Success		200	{array}		dto.DataResponse[models.Movie]
// @Failure		400	{object}	dto.ErrorResponse
// @Failure		404	{object}	dto.ErrorResponse
// @Router			/movies/{id} [get]
func GetMovieById(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, dto.ErrorResponse{
			BaseResponse: dto.BaseResponse{
				Message: "Invalid movie ID",
				Success: false,
			},
		})
		return
	}

	movie := database.FindMovieById(idInt)
	if movie != nil {
		c.IndentedJSON(http.StatusOK, dto.DataResponse[models.Movie]{
			BaseResponse: dto.BaseResponse{
				Message: "Movie found",
				Success: true,
			},
			Data: *movie,
		})
		return
	}

	c.IndentedJSON(http.StatusNotFound, dto.ErrorResponse{
		BaseResponse: dto.BaseResponse{
			Message: "Movie not found",
			Success: false,
		},
	})
}

// @Summary		Create a new movie
// @Description	Create a new movie
// @Tags			Movies
// @Param			movie	body		models.Movie	true	"Movie object to create"
// @Success		201		{object}	dto.DataResponse[models.Movie]
// @Failure		400		{object}	dto.ErrorResponse
// @Router			/movies [post]
func PostMovie(c *gin.Context) {
	var newMovie models.Movie

	if err := c.ShouldBindJSON(&newMovie); err != nil {
		c.IndentedJSON(http.StatusBadRequest, dto.ErrorResponse{
			BaseResponse: dto.BaseResponse{
				Message: "Invalid Movie body",
				Success: false,
			},
		})
		return
	}

	// assume primary key is the latest ID in the list
	database.FindMovieById(newMovie.Id)
	database.MovieId++
	newMovie.Id = database.MovieId

	database.Movies = append(database.Movies, newMovie)
	c.IndentedJSON(http.StatusCreated, dto.DataResponse[models.Movie]{
		BaseResponse: dto.BaseResponse{
			Message: "Movie created successfully",
			Success: true,
		},
		Data: newMovie,
	})
}

// @Summary		Update a movie
// @Description	Update a movie by ID
// @Tags			Movies
// @Param			id		path		int				true	"Movie ID"
// @Param			movie	body		models.Movie	true	"Updated movie object"
// @Success		200		{object}	dto.DataResponse[models.Movie]
// @Failure		400		{object}	dto.ErrorResponse
// @Failure		404		{object}	dto.ErrorResponse
// @Router			/movies/{id} [put]
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	var updatedMovie models.Movie

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, dto.ErrorResponse{
			BaseResponse: dto.BaseResponse{
				Message: "Invalid id",
				Success: false,
			},
		})
		return
	}

	if err := c.ShouldBindJSON(&updatedMovie); err != nil {
		c.IndentedJSON(http.StatusBadRequest, dto.ErrorResponse{
			BaseResponse: dto.BaseResponse{
				Message: "Invalid Movie body",
				Success: false,
			},
		})
		return
	}

	movie := database.FindMovieById(idInt)
	if movie == nil {
		c.IndentedJSON(http.StatusNotFound, dto.ErrorResponse{
			BaseResponse: dto.BaseResponse{
				Message: "Movie not found",
				Success: false,
			},
		})
		return
	}

	// check if id is updated, if so, check if it already exists
	if updatedMovie.Id != movie.Id {
		if database.FindMovieById(updatedMovie.Id) != nil {
			c.IndentedJSON(http.StatusBadRequest, dto.ErrorResponse{
				BaseResponse: dto.BaseResponse{
					Message: "Movie with updated ID already exists",
					Success: false,
				},
			})
			return
		}
	}

	*movie = updatedMovie
	c.IndentedJSON(http.StatusOK, dto.DataResponse[models.Movie]{
		BaseResponse: dto.BaseResponse{
			Message: "Movie updated successfully",
			Success: true,
		},
		Data: *movie,
	})
}

// @Summary		Delete a movie
// @Description	Delete a movie by ID
// @Tags			Movies
// @Param			id	path		int	true	"Movie ID"
// @Success		200	{object}	dto.BaseResponse
// @Failure		404	{object}	dto.ErrorResponse
// @Router			/movies/{id} [delete]
func DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	for i, movie := range database.Movies {
		if id == strconv.Itoa(movie.Id) {
			database.Movies = append(database.Movies[:i], database.Movies[i+1:]...)
			c.IndentedJSON(http.StatusOK, dto.ErrorResponse{
				BaseResponse: dto.BaseResponse{
					Message: "Movie deleted successfully",
					Success: false,
				},
			})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, dto.ErrorResponse{
		BaseResponse: dto.BaseResponse{
			Message: "Movie not found",
			Success: false,
		},
	})
}
