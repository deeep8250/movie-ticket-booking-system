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

func (r *TheaterRepository) GetShowsRepo(c context.Context, TheaterId int) ([]models.TheaterShows, error) {
	query := `select s.id as show_id,t.theater_name as theater_name,h.hall_name as hall_name,
        m.title as movie_name,m.language,s.starts_at,s.ends_at,s.base_price from theaters as t
		join halls as h on t.id=h.theater_id
		join shows as s on s.hall_id=h.id
		join movies as m on  s.movie_id=m.id
		where t.id=$1
		order by (t.id,h.id,m.id,s.id);
		`

	var shows []models.TheaterShows
	err := r.db.SelectContext(c, &shows, query, TheaterId)
	if err != nil {
		return nil, err
	}

	return shows, nil

}
