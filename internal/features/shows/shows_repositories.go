package shows

import (
	"context"
	"errors"

	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
	"github.com/jmoiron/sqlx"
)

type ShowRepositories struct {
	db *sqlx.DB
}

// theaters handlers
func NewShowsRepositories() *ShowRepositories {
	return &ShowRepositories{
		db: config.DBClients.PostgresClient,
	}
}

func (r *ShowRepositories) GetShowsRepo(c context.Context, TheaterId int) ([]models.TheaterShows, error) {

	// checking theater is existing or not for propper error handing
	var theaterCounter int
	q := `select count(*) from theaters where id=$1`
	err := r.db.GetContext(c, &theaterCounter, q, TheaterId)
	if err != nil {
		return nil, err
	} else if theaterCounter <= 0 {
		return nil, errors.New("theater not found")
	}

	query := `select s.id as show_id,t.theater_name as theater_name,h.hall_name as hall_name,
        m.title as movie_name,m.language,s.starts_at,s.ends_at,s.base_price from theaters as t
		join halls as h on t.id=h.theater_id
		join shows as s on s.hall_id=h.id
		join movies as m on  s.movie_id=m.id
		where t.id=$1
		order by (t.id,h.id,m.id,s.id);
		`

	var shows []models.TheaterShows
	err = r.db.SelectContext(c, &shows, query, TheaterId)
	if err != nil {
		return nil, err
	}

	return shows, nil

}
func (r *ShowRepositories) GetShowsByMovieIDRepo(c context.Context, MovieID int) ([]models.Shows, error) {

	query := `select sh.id AS show_id, t.theater_name,h.hall_name,t.city,sh.starts_at,sh.ends_at,sh.base_price from shows as sh
	          join halls as h on h.id=sh.hall_id 
			  join theaters as t on h.theater_id=t.id 
			  join movies as m on m.id=sh.movie_id 
			  where m.id=$1
			  `

	var AvlShows []models.Shows
	err := r.db.SelectContext(c, &AvlShows, query, MovieID)

	if err != nil {
		return nil, err
	}
	return AvlShows, nil
}
