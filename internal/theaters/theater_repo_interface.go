package theaters

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type TheaterRepoInterface interface {
	GetTheaters(c context.Context) ([]models.Theater, error)
	GetShowsRepo(c context.Context, TheaterId int) ([]models.TheaterShows, error)
	GetSeats(c context.Context, showsId int) (*models.SeatsInShows, error)
	BookSeat(c context.Context, userID, showID int, seats []int) (*models.SeatBooking, error)
	GetBookingByBookingsId(c context.Context, bookingID, userID int) (*models.BookingsDetails, error)
	UserBookingHistoryRepo(c context.Context, userID int) (string, []models.BookingHistory, error)
	BookingCancelRepo(c context.Context, BookingID, userID int) (*models.Bookings, error)
	GetMoviesRepo(c context.Context) ([]models.GetMovies, error)
	GetMovieByIDRepo(c context.Context, id int) (*models.GetMovies, error)
	GetShowsByMovieIDRepo(c context.Context, MovieID int) ([]models.Shows, error)

	SeatLockRepo(c context.Context, userID int, showID int, seats []int) error
}
