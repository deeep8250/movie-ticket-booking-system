package movies

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type MoviesServicesInterface interface {
	GetMoviesService(c context.Context) ([]dto.GetMovies, error)
	GetMovieByIDService(c context.Context, id int) (*dto.GetMovies, error)
}
