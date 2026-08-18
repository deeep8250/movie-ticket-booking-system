package shows

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type ShowsServicesInterface interface {
	GetShowsService(c context.Context, TheaterId int) ([]dto.TheaterShows, error)
	GetShowsByMovieIDService(c context.Context, MovieID int) ([]dto.Shows, error)
}
