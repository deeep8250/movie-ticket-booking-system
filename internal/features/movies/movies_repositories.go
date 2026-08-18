package movies

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
	"github.com/jmoiron/sqlx"
)

type MoviesRepositories struct {
	db *sqlx.DB
}

// theaters handlers
func NewMoviesRepositories() *MoviesRepositories {
	return &MoviesRepositories{
		db: config.DBClients.PostgresClient,
	}
}

// movies
func (r *MoviesRepositories) GetMoviesRepo(c context.Context) ([]models.GetMovies, error) {
	var movies []models.GetMovies
	query := `select id,title,language,duration_min,release_date from movies`
	err := r.db.SelectContext(c, &movies, query)
	if err != nil {
		return nil, err
	}

	return movies, err
}
func (r *MoviesRepositories) GetMovieByIDRepo(c context.Context, id int) (*models.GetMovies, error) {

	var movie models.GetMovies

	query := `select id,title,language,duration_min,release_date from movies where id=$1`
	err := r.db.GetContext(c, &movie, query, id)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}
