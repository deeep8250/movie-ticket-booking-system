package bookings

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type BookingServices struct {
	repo BookingRepoInterface
}

// theaters handlers
func NewBookingServices(S BookingRepoInterface) *BookingServices {
	return &BookingServices{
		repo: S,
	}
}

func (s *BookingServices) BookSeatService(c context.Context, userID, showID int, seats []int) (*dto.SeatBooking, error) {
	bookingData, err := s.repo.BookSeat(c, userID, showID, seats)
	if err != nil {
		return nil, err
	}

	bookingD := dto.SeatBooking{
		Id:          bookingData.Id,
		UserId:      bookingData.UserId,
		ShowID:      bookingData.ShowID,
		TotalPrice:  bookingData.TotalPrice,
		SeatsBooked: bookingData.SeatsBooked,
	}

	return &bookingD, nil
}

func (s *BookingServices) GetBookingByBookingsIdService(c context.Context, bookingID, userID int) (*dto.BookingsDetails, error) {

	bookingDetails, err := s.repo.GetBookingByBookingsId(c, bookingID, userID)
	if err != nil {
		return nil, err
	}

	bookingDetailsDTo := dto.BookingsDetails{

		BookingID:   bookingDetails.BookingID,
		UserID:      bookingDetails.UserID,
		UserName:    bookingDetails.UserName,
		ShowID:      bookingDetails.ShowID,
		MovieName:   bookingDetails.MovieName,
		TheaterName: bookingDetails.TheaterName,
		HallName:    bookingDetails.HallName,
		SeatID:      bookingDetails.SeatID,
		SeatNumber:  bookingDetails.SeatNumber,
		Status:      bookingDetails.Status,
		TotalAmount: int(bookingDetails.TotalAmount),
		CreatedAt:   bookingDetails.CreatedAt,
	}

	return &bookingDetailsDTo, nil

}
func (s *BookingServices) UserBookingHistoryService(c context.Context, userID int) (*dto.BookingHistoryWithUsrDtlsAtchd, error) {

	Username, BookingHistoriesFromRepo, err := s.repo.UserBookingHistoryRepo(c, userID)
	if err != nil {
		return nil, err
	}
	bookingHistoryResponse := dto.BookingHistoryWithUsrDtlsAtchd{
		UserID:   userID,
		Username: Username,
	}

	for _, history := range BookingHistoriesFromRepo {

		bookingHistoryDTO := dto.BookingHistory{
			BookingID:   history.BookingID,
			CreatedAt:   history.CreatedAt,
			HallName:    history.HallName,
			MovieName:   history.MovieName,
			SeatCount:   history.SeatCount,
			Status:      history.Status,
			TheaterName: history.TheaterName,
			TotalAmount: history.TotalAmount,
		}
		bookingHistoryResponse.BookingHitories = append(bookingHistoryResponse.BookingHitories, bookingHistoryDTO)

	}

	return &bookingHistoryResponse, nil

}
func (s *BookingServices) BookingCancelService(c context.Context, BookingID, userID int) (*dto.BookingsDto, error) {
	bookingCancelationDetails, err := s.repo.BookingCancelRepo(c, BookingID, userID)
	if err != nil {
		return nil, err
	}
	bd := &dto.BookingsDto{
		BookingID:   bookingCancelationDetails.BookingID,
		UserID:      bookingCancelationDetails.UserID,
		ShowID:      bookingCancelationDetails.ShowID,
		Status:      bookingCancelationDetails.Status,
		TotalAmount: bookingCancelationDetails.TotalAmount,
		CreatedAt:   bookingCancelationDetails.CreatedAt,
		UpdatedAt:   bookingCancelationDetails.UpdatedAt,
	}
	return bd, nil
}
