package routes

import (
	"github.com/deeep8250/movie-ticket-booking-system/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.GET("/health", handlers.CheckHealth)
	r.GET("/ready", handlers.CheckReady)
}
