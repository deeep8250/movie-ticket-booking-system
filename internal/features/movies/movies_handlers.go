package movies

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type moviesHandlers struct {
	service MoviesServicesInterface
}

// theaters handlers
func NewMoviesHandlers(S MoviesServicesInterface) *moviesHandlers {
	return &moviesHandlers{
		service: S,
	}
}

func (h *moviesHandlers) GetMoviesHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	movies, err := h.service.GetMoviesService(ctx)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"movies_available": movies,
	})

}
func (h *moviesHandlers) GetMoviesByIDHandler(c *gin.Context) {
	MovieID := c.Param("id")
	MovieIDint, err := strconv.Atoi(MovieID)
	if err != nil || MovieIDint <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invlid input",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	movie, err := h.service.GetMovieByIDService(ctx, MovieIDint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
}
