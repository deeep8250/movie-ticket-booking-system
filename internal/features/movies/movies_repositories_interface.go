package movies

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type MoviesRepositoriesInterface interface {
	GetMoviesRepo(c context.Context) ([]models.GetMovies, error)
	GetMovieByIDRepo(c context.Context, id int) (*models.GetMovies, error)
}
