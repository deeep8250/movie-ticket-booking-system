package bookings

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type BookingServiceInterface interface {
	BookSeatService(c context.Context, userID, showID int, seats []int) (*dto.SeatBooking, error)
	GetBookingByBookingsIdService(c context.Context, bookingID, userID int) (*dto.BookingsDetails, error)
	UserBookingHistoryService(c context.Context, userID int) (*dto.BookingHistoryWithUsrDtlsAtchd, error)
	BookingCancelService(c context.Context, BookingID, useID int) (*dto.BookingsDto, error)
}
