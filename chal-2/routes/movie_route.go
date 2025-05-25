package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sglkc/roketin-be-test/chal-2/controllers"
)

func RegisterMovieRoutes(router *gin.Engine) {
	router.GET("/movies", controllers.GetMovies)
	router.GET("/movies/:id", controllers.GetMovieById)
	router.GET("/movies/search", controllers.SearchMovie)
	router.POST("/movies", controllers.PostMovie)
	router.PUT("/movies/:id", controllers.UpdateMovie)
	router.DELETE("/movies/:id", controllers.DeleteMovie)
}
