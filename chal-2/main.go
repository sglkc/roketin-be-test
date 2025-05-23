/* https://go.dev/doc/tutorial/web-service-gin
 * • API to create and upload movies. Required information related with a movies are at least
 *   title, description, duration, artists, genres
 * • API to update movie
 * • API to list all movies with pagination
 * • API to search movie by title/description/artists/genres
 */
package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Duration    int      `json:"duration"` // in minutes
	Artists     []string `json:"artists"`
	Genres      []string `json:"genres"`
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

var id int = 0

func paginate(c *gin.Context, movies []Movie) []Movie {
	page := c.Query("page")
	limit := c.Query("limit")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

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

func getMovies(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, paginate(c, movies))
}

func getMovieById(c *gin.Context) {
	id := c.Param("id")

	for _, movie := range movies {
		if id == strconv.Itoa(movie.Id) {
			c.IndentedJSON(http.StatusOK, movie)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}

func postMovie(c *gin.Context) {
	var newMovie Movie

	if err := c.BindJSON(&newMovie); err != nil {
		return
	}

	id++
	newMovie.Id = id

	movies = append(movies, newMovie)
	c.IndentedJSON(http.StatusCreated, newMovie)
}

func updateMovie(c *gin.Context) {
	id := c.Param("id")
	var updatedMovie Movie

	if err := c.BindJSON(&updatedMovie); err != nil {
		return
	}

	for i, movie := range movies {
		if id == strconv.Itoa(movie.Id) {
			movies[i] = updatedMovie
			c.IndentedJSON(http.StatusOK, updatedMovie)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "movie not found"})
}

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
		if movie.Id > id {
			id = movie.Id
		}
	}
	router := gin.Default()
	router.GET("/movies", getMovies)
	router.GET("/movies/:id", getMovieById)
	router.GET("/movies/search", searchMovie)
	router.POST("/movies", postMovie)
	router.PUT("/movies/:id", updateMovie)
	router.DELETE("/movies/:id", deleteMovie)

	router.Run("localhost:8080")
}
