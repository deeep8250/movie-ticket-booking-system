package theaters

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type TheaterService struct {
	repo TheaterRepoInterface
}

func NewTheaterService(Repo TheaterRepoInterface) *TheaterService {
	return &TheaterService{
		repo: Repo,
	}
}

func (s *TheaterService) GetTheatersService(c context.Context) ([]dto.Theater, error) {
	theatersValues, err := s.repo.GetTheaters(c)
	if err != nil {
		return nil, err
	}

	var theaters []dto.Theater
	for r := range theatersValues {

		theater := dto.Theater{
			Id:            theatersValues[r].Id,
			TheaterName:   theatersValues[r].TheaterName,
			NumberOfHalls: theatersValues[r].NumberOfHalls,
			City:          theatersValues[r].City,
			State:         theatersValues[r].State,
		}
		theaters = append(theaters, theater)

	}

	return theaters, nil
}

func (s *TheaterService) GetShowsService(c context.Context, TheaterId int) ([]dto.TheaterShows, error) {
	shows, err := s.repo.GetShowsRepo(c, TheaterId)
	if err != nil {
		return nil, err
	}

	var showsInTheater []dto.TheaterShows
	for r := range shows {

		showInTheater := dto.TheaterShows{
			ShowId:      shows[r].ShowId,
			TheaterName: shows[r].TheaterName,
			HallName:    shows[r].HallName,
			MovieName:   shows[r].MovieName,
			Language:    shows[r].Language,
			StartDate:   shows[r].StartDate,
			EndDate:     shows[r].EndDate,
			Price:       shows[r].Price,
		}
		showsInTheater = append(showsInTheater, showInTheater)

	}

	return showsInTheater, nil
}

func (s *TheaterService) GetSeatsService(c context.Context, showsId int) (*dto.SeatsInShows, error) {
	seatsAvailableForTheShow, err := s.repo.GetSeats(c, showsId)
	if err != nil {
		return nil, err
	}

	seats := dto.SeatsInShows{
		ShowId:    seatsAvailableForTheShow.ShowId,
		MovieName: seatsAvailableForTheShow.MovieName,
		HallName:  seatsAvailableForTheShow.HallName,
	}

	for r := range seatsAvailableForTheShow.SeatsAvailable {
		seats.SeatsAvailable = append(seats.SeatsAvailable, dto.Seats(seatsAvailableForTheShow.SeatsAvailable[r]))

	}

	return &seats, nil
}

func (s *TheaterService) BookSeatService(c context.Context, userID, showID int, seats []int) (*dto.SeatBooking, error) {
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

func (s *TheaterService) GetBookingByBookingsIdService(c context.Context, bookingID, userID int) (*dto.BookingsDetails, error) {

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

func (s *TheaterService) UserBookingHistoryService(c context.Context, userID int) (*dto.BookingHistoryWithUsrDtlsAtchd, error) {

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

func (s *TheaterService) BookingCancelService(c context.Context, BookingID, userID int) (*dto.BookingsDto, error) {
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

func (s *TheaterService) GetMoviesService(c context.Context) ([]dto.GetMovies, error) {
	movies, err := s.repo.GetMoviesRepo(c)
	if err != nil {
		return nil, err
	}

	var ReturnedMovies []dto.GetMovies
	for _, movie := range movies {
		m := dto.GetMovies{
			MoviesID:    movie.MoviesID,
			Title:       movie.Title,
			Language:    movie.Language,
			DurationMin: movie.DurationMin,
			ReleaseDate: movie.ReleaseDate,
		}

		ReturnedMovies = append(ReturnedMovies, m)
	}

	return ReturnedMovies, nil

}

func (s *TheaterService) GetMovieByIDService(c context.Context, id int) (*dto.GetMovies, error) {

	m, err := s.repo.GetMovieByIDRepo(c, id)
	if err != nil {
		return nil, err
	}
	movie := dto.GetMovies{
		MoviesID:    m.MoviesID,
		Title:       m.Title,
		Language:    m.Language,
		DurationMin: m.DurationMin,
		ReleaseDate: m.ReleaseDate,
	}

	return &movie, nil
}

func (s *TheaterService) GetShowsByMovieIDService(c context.Context, MovieID int) ([]dto.Shows, error) {

	sh, err := s.repo.GetShowsByMovieIDRepo(c, MovieID)
	if err != nil {
		return nil, err
	}

	var shows []dto.Shows
	for _, show := range sh {
		s := dto.Shows{
			ShowID:      show.ShowID,
			TheaterName: show.TheaterName,
			City:        show.City,
			HallName:    show.HallName,
			StartsAt:    show.StartsAt,
			EndsAt:      show.EndsAt,
			BasePrice:   show.BasePrice,
		}
		shows = append(shows, s)
	}

	return shows, nil

}

func (s *TheaterService) SeatLockService(c context.Context, userID int, showID int, seats []int) error {
	err := s.repo.SeatLockRepo(c, userID, showID, seats)
	if err != nil {
		return err
	}
	return nil
}
func (s *TheaterService) SeatUnLockService(c context.Context, userID int, showID int, seatIDs []int) error {
	err := s.repo.SeatUnLockRepo(c, userID, showID, seatIDs)
	if err != nil {
		return err
	}
	return nil
}
