package routes

import (
	"github.com/deeep8250/movie-ticket-booking-system/internal/auth"
	"github.com/deeep8250/movie-ticket-booking-system/internal/features/bookings"
	"github.com/deeep8250/movie-ticket-booking-system/internal/features/movies"
	"github.com/deeep8250/movie-ticket-booking-system/internal/features/seats"
	"github.com/deeep8250/movie-ticket-booking-system/internal/features/shows"
	"github.com/deeep8250/movie-ticket-booking-system/internal/features/theaters"
	"github.com/deeep8250/movie-ticket-booking-system/internal/middleware"
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

	// movies DI
	movieRepo := movies.NewMoviesRepositories()
	movieServices := movies.NewMoviesServices(movieRepo)
	movieHandlers := movies.NewMoviesHandlers(movieServices)

	// shows DI
	showRepo := shows.NewShowsRepositories()
	showServices := shows.NewShowsServices(showRepo)
	showHandlers := shows.NewShowsHandlers(showServices)

	//seats DI
	seatRepo := seats.NewSeatsRepositories()
	seatServices := seats.NewSeatsServices(seatRepo)
	seatHandlers := seats.NewSeatsHandlers(seatServices)

	// booking DI
	bookingRepo := bookings.NewBookingRepo()
	bookingServices := bookings.NewBookingServices(bookingRepo)
	bookingHandlers := bookings.NewBookingsHandlers(bookingServices)

	pb := r.Group("/public")
	{

		//public
		// movies
		pb.GET("/movies", movieHandlers.GetMoviesHandler)
		pb.GET("/movies/:id", movieHandlers.GetMoviesByIDHandler)

		// shows
		pb.GET("/movies/:id/shows", showHandlers.GetShowByMovieIdHandler)
		pb.GET("/theaters/shows/:id", showHandlers.GetShows)

		// theaters
		pb.GET("/theaters", theaterHandler.GetTheaters)

		//seats
		pb.GET("/theaters/shows/:id/seats", seatHandlers.GetSeatsHandler)
		pb.POST("/signup", authHandler.CreateUserHandler)
		pb.POST("/login", authHandler.LoginHandler)

	}

	pv := r.Group("/private", middleware.Middleware())
	{
		//bokings
		pv.POST("/bookings", bookingHandlers.BookSeatHandler)
		pv.GET("/bookings/:id/details", bookingHandlers.GetBookingDetailsFromId)
		pv.GET("/users/me/bookings", bookingHandlers.UserBookingHistory)
		pv.PATCH("/bookings/:id/cancel", bookingHandlers.BookingCancelation)

		//auth
		pv.GET("/user", authHandler.GetUserHandler)

		// seats
		pv.POST("/shows/:id/seats/lock", seatHandlers.SeatLockHandler)
		pv.POST("/shows/:id/seats/unlock", seatHandlers.SeatUnLockHandler)

	}

}
