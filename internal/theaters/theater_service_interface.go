package theaters

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type TheaterServiceInterface interface {
	GetTheatersService(c context.Context) ([]dto.Theater, error)
	GetShowsService(c context.Context, TheaterId int) ([]dto.TheaterShows, error)
	GetSeatsService(c context.Context, showsId int) (*dto.SeatsInShows, error)
	BookSeatService(c context.Context, userID, showID int, seats []int) (*dto.SeatBooking, error)
	GetBookingByBookingsIdService(c context.Context, bookingID int) (*dto.BookingsDetails, error)
	UserBookingHistoryService(c context.Context, userID int) (*dto.BookingHistoryWithUsrDtlsAtchd, error)
	BookingCancelService(c context.Context, userID int) (*dto.BookingsDto, error)
	GetMoviesService(c context.Context) ([]dto.GetMovies, error)
}
