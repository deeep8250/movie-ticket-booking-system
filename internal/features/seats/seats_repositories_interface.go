package seats

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type SeatRepositoriesInterface interface {
	GetSeats(c context.Context, showsId int) (*models.SeatsInShows, error)
	SeatLockRepo(c context.Context, userID int, showID int, seatIDs []int) error
	SeatUnLockRepo(c context.Context, userID int, showID int, seatIDs []int) error
}
