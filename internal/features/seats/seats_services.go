package seats

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type SeatsServices struct {
	repo SeatRepositoriesInterface
}

// theaters handlers
func NewSeatsServices(r SeatRepositoriesInterface) *SeatsServices {
	return &SeatsServices{
		repo: r,
	}
}

func (s *SeatsServices) GetSeatsService(c context.Context, showsId int) (*dto.SeatsInShows, error) {
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

func (s *SeatsServices) SeatLockService(c context.Context, userID int, showID int, seats []int) error {
	err := s.repo.SeatLockRepo(c, userID, showID, seats)
	if err != nil {
		return err
	}
	return nil
}
func (s *SeatsServices) SeatUnLockService(c context.Context, userID int, showID int, seatIDs []int) error {
	err := s.repo.SeatUnLockRepo(c, userID, showID, seatIDs)
	if err != nil {
		return err
	}
	return nil
}
