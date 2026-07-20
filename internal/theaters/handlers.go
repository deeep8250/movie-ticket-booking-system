package theaters

import (
	"context"
	"net/http"
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
