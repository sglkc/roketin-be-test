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

	"github.com/gin-gonic/gin"
	"github.com/sglkc/roketin-be-test/chal-2/routes"
)

// @title			Movies API
// @version		1.0
// @description	This is a sample movies API using Gin framework.
// @contact.name	sglkc
// @contact.url	https://github.com/sglkc/roketin-be-test
// @produce		json
// @accept			json
func main() {
	router := gin.Default()

	routes.RegisterSwaggerRoutes(router)
	routes.RegisterMovieRoutes(router)

	log.Println("Running at localhost:8080 (docs at http://localhost:8080/swagger/index.html)")
	router.Run("localhost:8080")
}
