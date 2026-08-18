package shows

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type ShowsRepositoriesInterface interface {
	GetShowsRepo(c context.Context, TheaterId int) ([]models.TheaterShows, error)
	GetShowsByMovieIDRepo(c context.Context, MovieID int) ([]models.Shows, error)
}
