package theaters

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type TheaterServiceInterface interface {
	GetTheatersService(c context.Context) ([]dto.Theater, error)
	GetShowsService(c context.Context, TheaterId int) ([]dto.TheaterShows, error)
	GetSeatsService(c context.Context, showsId int) (*models.SeatsInShows, error)
}
