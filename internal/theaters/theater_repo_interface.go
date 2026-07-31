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
	GetBookingByBookingsId(c context.Context, bookingID int) (*models.BookingsDetails, error)
	UserBookingHistoryRepo(c context.Context, userID int) (string, []models.BookingHistory, error)
	BookingCancelRepo(c context.Context, BookingID int) (*models.Bookings, error)
}
