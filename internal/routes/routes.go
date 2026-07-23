package routes

import (
	"github.com/deeep8250/movie-ticket-booking-system/internal/theaters"
	statusandhealth "github.com/deeep8250/movie-ticket-booking-system/status_and_health"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.GET("/health", statusandhealth.CheckHealth)
	r.GET("/ready", statusandhealth.CheckReady)

	theaterRepo := theaters.NewTheaterRepo()
	theaterService := theaters.NewTheaterService(theaterRepo)
	theaterHandler := theaters.NewTheaterHandler(theaterService)

	p := r.Group("/public")
	{
		p.GET("/theaters", theaterHandler.GetTheaters)
		p.GET("/theaters/shows/:id", theaterHandler.GetShows)
		p.GET("/theaters/shows/:id/seats", theaterHandler.GetSeatsHandler)
		p.POST("/bookings", theaterHandler.BookSeatHandler)
	}
}
