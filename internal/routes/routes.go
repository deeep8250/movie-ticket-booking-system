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

	pb := r.Group("/public")
	{

		//public
		pb.GET("/movies", theaterHandler.GetMoviesHandler)
		pb.GET("/movies/:id", theaterHandler.GetMoviesByIDHandler)
		pb.GET("/movies/:id/shows", theaterHandler.GetShowByMovieIdHandler)
		pb.GET("/theaters", theaterHandler.GetTheaters)
		pb.GET("/theaters/shows/:id", theaterHandler.GetShows)
		pb.GET("/theaters/shows/:id/seats", theaterHandler.GetSeatsHandler)
		pb.POST("/signup", authHandler.CreateUserHandler)
		pb.POST("/login", authHandler.LoginHandler)

	}

	pv := r.Group("/private", middleware.Middleware())
	{
		pv.POST("/bookings", theaterHandler.BookSeatHandler)
		pv.GET("/bookings/:id/details", theaterHandler.GetBookingDetailsFromId)
		pv.GET("/users/me/bookings", theaterHandler.UserBookingHistory)
		pv.PATCH("/bookings/:id/cancel", theaterHandler.BookingCancelation)
		pv.GET("/user", authHandler.GetUserHandler)
		pv.POST("/shows/:id/seats/lock", theaterHandler.SeatLockHandler)
		pv.POST("/shows/:id/seats/unlock", theaterHandler.SeatUnLockHandler)

	}

}
