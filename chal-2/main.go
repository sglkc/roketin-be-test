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

	"github.com/gin-gonic/gin"
)

type Movie struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Duration    int      `json:"duration"` // in minutes
	Artists     []string `json:"artists"`
	Genres      []string `json:"genres"`
}

var movies = []Movie{
	{
		Title:       "Final Destination: Bloodlines",
		Description: "Plagued by a recurring violent nightmare, a college student returns home to find the one person who can break the cycle and save her family from the horrific fate that inevitably awaits them.",
		Duration:    90,
		Artists:     []string{"Kaitlyn Santa Juana", "Teo Briones", "Rya Kihlstedt"},
		Genres:      []string{"Horror", "Splatter Horror"},
	},
	{
		Title:       "Mission: Impossible - The Final Reckoning",
		Description: "Our lives are the sum of our choices. Tom Cruise is Ethan Hunt in Mission: Impossible - The Final Reckoning.",
		Duration:    169,
		Artists:     []string{"Tom Cruise", "Haylett Atwell", "Ving Rhames"},
		Genres:      []string{"Action", "Adeventure", "Thriller"},
	},
}

func getMovies(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, movies)
}

func main() {
	router := gin.Default()
	router.GET("/movies", getMovies)

	router.Run("localhost:8080")
}
