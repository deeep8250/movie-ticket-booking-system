package routes

import (
	"github.com/deeep8250/movie-ticket-booking-system/internal/auth"
	"github.com/deeep8250/movie-ticket-booking-system/internal/middleware"
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

	authRepo := auth.NewAuthRepo()
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)
	r.POST("/user/signup", authHandler.CreateUserHandler)
	r.POST("/login", authHandler.LoginHandler)
	r.GET("/user", middleware.Middleware(), authHandler.GetUserHandler)
	p := r.Group("/public", middleware.Middleware())
	{

		p.GET("/theaters", theaterHandler.GetTheaters)
		p.GET("/theaters/shows/:id", theaterHandler.GetShows)
		p.GET("/theaters/shows/:id/seats", theaterHandler.GetSeatsHandler)
		p.POST("/bookings", theaterHandler.BookSeatHandler)
		p.GET("/bookings/:id/details", theaterHandler.GetBookingDetailsFromId)
		p.GET("/users/:id/bookings", theaterHandler.UserBookingHistory)
		p.PATCH("/bookings/:id/cancel", theaterHandler.BookingCancelation)
		p.GET("/movies", theaterHandler.GetMoviesHandler)
		p.GET("/movies/:id", theaterHandler.GetMoviesByIDHandler)
		p.GET("/movies/:id/shows", theaterHandler.GetShowByMovieIdHandler)
	}
}
