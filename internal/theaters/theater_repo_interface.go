package theaters

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type TheaterRepoInterface interface {
	GetTheaters(c context.Context) ([]models.Theater, error)
	GetShowsRepo(c context.Context, TheaterId int) ([]models.TheaterShows, error)
}
