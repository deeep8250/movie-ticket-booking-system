package theaters

import (
	"context"
	"errors"

	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (r *TheaterRepository) GetSeats(c context.Context, showsId int) (*models.SeatsInShows, error) {

	var Seats []models.Seats
	query := `SELECT 
		s.id AS seat_id,
		s.seat_number,
		s.seat_type,
		CASE
			WHEN s.is_active = false THEN 'disabled'
			WHEN sb.id IS NOT NULL THEN 'booked'
			ELSE 'available'
		END AS status


	FROM seats AS s
	JOIN shows AS sh 
		ON sh.hall_id = s.hall_id


	LEFT JOIN seat_bookings AS sb 
	   ON sb.seat_id = s.id
		AND sb.show_id = sh.id


		
	WHERE sh.id = $1
	ORDER BY s.id;`

	err := r.db.SelectContext(c, &Seats, query, showsId)
	if err != nil {
		return nil, err
	}
	// showId,movie_name,hall_name
	query2 := `select s.id as show_id,m.title as movie_name,h.hall_name from halls as h
	            join shows as s on h.id=s.hall_id
				join movies as m on s.movie_id=m.id where s.id=$1 order by (h.id,s.id,m.id)  `

	var SeatinShows models.SeatsInShows
	err = r.db.GetContext(c, &SeatinShows, query2, showsId)
	if err != nil {
		return nil, err
	}
	SeatinShows.SeatsAvailable = Seats
	return &SeatinShows, nil

}

func (r *TheaterRepository) BookSeat(c context.Context, userID, showID int, seats []int) (*models.SeatBooking, error) {

	tx, err := r.db.BeginTxx(c, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// checking if the given user id is valid or not
	var count int
	err = tx.GetContext(c, &count, `select count(*) from users where id=$1`, userID)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("user not found")
	}

	// 1 validate seats belongs to the show and active also
	var validSeatCount int
	err = tx.GetContext(c, &validSeatCount, `select count(*) from seats as s join shows as sh on sh.hall_id=s.hall_id where sh.id=$1 and s.id=any($2) and s.is_active=true`, showID, pq.Array(seats))
	if err != nil {
		return nil, err
	}

	if validSeatCount != len(seats) {
		return nil, errors.New("one or more seats are invalid for this show")
	}

	// 2 check already booked seats

	var bookedSeatIds []int
	err = tx.SelectContext(c, &bookedSeatIds, `select seat_id from seat_bookings where seat_id=any($1) and show_id=$2`, pq.Array(seats), showID)
	if err != nil {
		return nil, err
	}

	if len(bookedSeatIds) > 0 {
		return nil, errors.New("one or more seats are already booked")
	}

	// 3. Get show price
	var basePrice float64
	err = tx.GetContext(c, &basePrice, `select base_price from shows where id=$1`, showID)
	if err != nil {
		return nil, err
	}

	totalAmount := float64(len(seats)) * basePrice

	// 4. Create booking row

	var bookingID int
	err = tx.QueryRowxContext(c, `insert into bookings (user_id,show_id,status,total_amount) values($1,$2,'confirmed',$3) returning id`, userID, showID, totalAmount).Scan(&bookingID)
	if err != nil {
		return nil, err
	}

	// 5. Insert selected seats

	for _, seatID := range seats {
		_, err := tx.ExecContext(c, `insert into seat_bookings(booking_id,show_id,seat_id) values($1,$2,$3)`, bookingID, showID, seatID)
		if err != nil {
			return nil, err
		}

	}

	// 5.1 fetching the final data that need to return
	var data []int

	err = tx.SelectContext(c, &data, `select s.id as seat_booked from seat_bookings as s where booking_id=$1`, bookingID)
	if err != nil {
		return nil, err
	}

	var bookingData models.SeatBooking
	// user_id,show_id,total_price
	err = tx.GetContext(c, &bookingData, `select id,user_id,show_id,total_amount from bookings where id=$1`, bookingID)
	if err != nil {
		return nil, err
	}

	bookingData.SeatsBooked = append(bookingData.SeatsBooked, data...)

	// 6. final save
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &bookingData, err

}
