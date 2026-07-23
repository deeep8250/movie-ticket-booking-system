package theaters

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TheaterHandler struct {
	service TheaterServiceInterface
}

func NewTheaterHandler(S TheaterServiceInterface) *TheaterHandler {
	return &TheaterHandler{
		service: S,
	}
}

func (h *TheaterHandler) GetTheaters(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()
	theaters, err := h.service.GetTheatersService(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"theaters": theaters,
	})

}

func (h *TheaterHandler) GetShows(c *gin.Context) {

	TheaterId := c.Param("id")
	theaterIdInt, err := strconv.Atoi(TheaterId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "theater id is invalid",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()

	shows, err := h.service.GetShowsService(ctx, theaterIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"TheaterID":       theaterIdInt,
		"shows_available": shows,
	})
}

func (h *TheaterHandler) GetSeatsHandler(c *gin.Context) {

	ShowId := c.Param("id")
	ShowIdInt, err := strconv.Atoi(ShowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid show id",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer cancel()
	SeatsAvailableForShows, err := h.service.GetSeatsService(ctx, ShowIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SeatsAvailableForShows)

}
