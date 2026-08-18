package shows

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ShowHandlers struct {
	service ShowsServicesInterface
}

// theaters handlers
func NewShowsHandlers(r ShowsServicesInterface) *ShowHandlers {
	return &ShowHandlers{
		service: r,
	}
}

func (h *ShowHandlers) GetShows(c *gin.Context) {

	TheaterId := c.Param("id")
	theaterIdInt, err := strconv.Atoi(TheaterId)
	if err != nil || theaterIdInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid theater id",
		})
		return

	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	shows, err := h.service.GetShowsService(ctx, theaterIdInt)
	if err != nil {

		errMsg := err.Error()
		if strings.Contains(errMsg, "theater not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": errMsg,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return

	}

	if len(shows) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"theater_id": theaterIdInt,
			"message":    "no shows available right at the moment",
			"shows":      shows,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"TheaterID":       theaterIdInt,
			"shows_available": shows,
		})
	}

}
func (h *ShowHandlers) GetShowByMovieIdHandler(c *gin.Context) {
	movieID := c.Param("id")
	movieIDInt, err := strconv.Atoi(movieID)
	if err != nil || movieIDInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid input",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	shows, err := h.service.GetShowsByMovieIDService(ctx, movieIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"eror": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie_id": movieIDInt,
		"shows":    shows,
	})

}
