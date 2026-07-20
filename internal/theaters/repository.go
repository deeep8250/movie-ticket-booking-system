package theaters

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
	"github.com/jmoiron/sqlx"
)

type TheaterRepository struct {
	db *sqlx.DB
}

func NewTheaterRepo() *TheaterRepository {
	return &TheaterRepository{
		db: config.DBClients.PostgresClient,
	}
}

func (r *TheaterRepository) GetTheaters(c context.Context) ([]models.Theater, error) {
	var theaterList []models.Theater
	query := `
		SELECT 
			t.id,
			t.theater_name,
			t.city,
			t.state,
			COUNT(h.id) AS number_of_halls
		FROM theaters t
		LEFT JOIN halls h ON h.theater_id = t.id
		GROUP BY t.id, t.theater_name, t.city, t.state
		ORDER BY t.id;
	`
	err := r.db.SelectContext(c, &theaterList, query)
	if err != nil {
		return nil, err
	}
	return theaterList, nil
}
