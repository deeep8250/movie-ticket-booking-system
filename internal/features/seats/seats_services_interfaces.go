package seats

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type SeatServicesInterface interface {
	GetSeatsService(c context.Context, showsId int) (*dto.SeatsInShows, error)
	SeatLockService(c context.Context, userID int, showID int, seats []int) error
	SeatUnLockService(c context.Context, userID int, showID int, seatIDs []int) error
}
