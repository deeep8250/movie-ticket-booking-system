package routes

import (
	"github.com/deeep8250/movie-ticket-booking-system/internal/health"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.GET("/health", health.CheckHealth)
}
