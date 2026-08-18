package bookings

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type BookingRepoInterface interface {
	BookSeat(c context.Context, userID, showID int, seats []int) (*models.SeatBooking, error)
	GetBookingByBookingsId(c context.Context, bookingID, userID int) (*models.BookingsDetails, error)
	UserBookingHistoryRepo(c context.Context, userID int) (string, []models.BookingHistory, error)
	BookingCancelRepo(c context.Context, BookingID, userID int) (*models.Bookings, error)
}
