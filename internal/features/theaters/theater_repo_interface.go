package theaters

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type TheaterRepoInterface interface {
	//theaters
	GetTheatersRepo(c context.Context) ([]models.Theater, error)
}
