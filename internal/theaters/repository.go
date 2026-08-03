package theaters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

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

func (r *TheaterRepository) GetSeats(c context.Context, showsId int) (*models.SeatsInShows, error) {

	//  checking if the given show id is valid or not
	var ShowCount int
	q := `select count(*) from shows where id=$1`
	err := r.db.GetContext(c, &ShowCount, q, showsId)
	if err != nil {
		return nil, err
	} else if ShowCount <= 0 {
		return nil, errors.New("show not found")
	}

	//

	var Seats []models.Seats
	// query := `SELECT
	// 	s.id AS seat_id,
	// 	s.seat_number,
	// 	s.seat_type,
	// 	CASE
	// 		WHEN s.is_active = false THEN 'disabled'
	// 		WHEN sb.id IS NOT NULL   AND b.status='confirmed' THEN 'booked'
	// 		WHEN sb.id IS NOT NULL   AND b.status='cancelled' THEN 'available'

	// 		ELSE 'available'
	// 	END AS status

	// FROM seats AS s
	// JOIN shows AS sh
	// 	ON sh.hall_id = s.hall_id

	// LEFT JOIN seat_bookings AS sb
	//    ON sb.seat_id = s.id
	// 	AND sb.show_id = sh.id

	// LEFT JOIN bookings as b
	//     ON b.id=sb.booking_id

	// WHERE sh.id = $1
	//  ORDER BY s.id;`

	query := `SELECT 
    s.id AS seat_id,
    s.seat_number,
    s.seat_type,
    CASE
        WHEN s.is_active = false THEN 'disabled'
        WHEN active_booking.seat_id IS NOT NULL THEN 'booked'
        ELSE 'available'
    END AS status
FROM seats AS s
JOIN shows AS sh 
    ON sh.hall_id = s.hall_id
LEFT JOIN (
    SELECT DISTINCT
        sb.seat_id,
        sb.show_id
    FROM seat_bookings AS sb
    JOIN bookings AS b
        ON b.id = sb.booking_id
    WHERE b.status IN ('confirmed', 'pending')
) AS active_booking
    ON active_booking.seat_id = s.id
    AND active_booking.show_id = sh.id
WHERE sh.id = $1
ORDER BY s.id;`

	err = r.db.SelectContext(c, &Seats, query, showsId)
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

	// checking if the showId is exists or not
	var ShowCount int
	err = tx.GetContext(c, &ShowCount, `select count(*) from shows where id=$1`, showID)
	if err != nil {
		return nil, err
	}
	if ShowCount == 0 {
		return nil, errors.New("show not found")
	}

	// 1 validate seats belongs to the show and active also
	var validSeatIDs []int
	err = tx.SelectContext(c, &validSeatIDs, `select s.id from seats as s join shows as sh on sh.hall_id=s.hall_id where sh.id=$1 and s.id=any($2) and s.is_active=true`, showID, pq.Array(seats))
	if err != nil {
		return nil, err
	}

	var invalidSeatIDs []int
	for i := range seats {
		NotPresent := !slices.Contains(validSeatIDs, seats[i])
		if NotPresent {
			invalidSeatIDs = append(invalidSeatIDs, seats[i])
		}
	}

	if len(invalidSeatIDs) > 0 {
		return nil, fmt.Errorf("invalid seats for this show: %v", invalidSeatIDs)
	}
	// 2 check already booked seats

	var bookedSeatIds []int
	err = tx.SelectContext(c, &bookedSeatIds,
		`SELECT sb.seat_id
         FROM seat_bookings AS sb
         JOIN bookings AS b
         ON b.id = sb.booking_id
         WHERE sb.seat_id = ANY($1)
         AND sb.show_id = $2
         AND b.status IN ('confirmed', 'pending');`, pq.Array(seats), showID)
	if err != nil {
		return nil, err
	}

	if len(bookedSeatIds) > 0 {
		return nil, fmt.Errorf("seats already booked: %v", bookedSeatIds)
	}

	// 3. Get show price
	var basePrice float64
	err = tx.GetContext(c, &basePrice, `select base_price from shows where id=$1`, showID)
	if err != nil {
		return nil, err
	}

	totalAmount := float64(len(seats)) * basePrice
	if totalAmount <= 10 {
		return nil, errors.New("invalid total amount")
	}

	// 4. Create booking row

	var bookingID int
	err = tx.QueryRowxContext(c, `insert into bookings (user_id,show_id,status,total_amount) values($1,$2,'confirmed',$3) returning id`, userID, showID, totalAmount).Scan(&bookingID)
	if err != nil {
		return nil, err
	}

	// 5. Insert selected seats

	for _, seatID := range seats {
		rows, err := tx.ExecContext(c, `insert into seat_bookings(booking_id,show_id,seat_id) values($1,$2,$3)`, bookingID, showID, seatID)
		if err != nil {
			return nil, err
		}

		rowsAffected, err := rows.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rowsAffected == 0 {
			return nil, fmt.Errorf("unable to book the seat %v", seatID)
		}

	}

	// 5.1 fetching the final data that need to return
	var data []int

	err = tx.SelectContext(c, &data, `select s.seat_id as seat_booked from seat_bookings as s where booking_id=$1`, bookingID)
	if err != nil {
		return nil, err
	}

	var bookingData models.SeatBooking
	//preparing return data
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

func (r *TheaterRepository) GetBookingByBookingsId(c context.Context, bookingID int) (*models.BookingsDetails, error) {
	var bookingDetailsCount int
	q1 := `select count(*) from bookings where id=$1`
	err := r.db.GetContext(c, &bookingDetailsCount, q1, bookingID)
	if err != nil {
		return nil, err
	}
	if bookingDetailsCount == 0 {
		return nil, errors.New("booking not found")
	}

	// fetching seat id and seat numbers according to the booking id
	type SeatDetails struct {
		SeatIds     int    `db:"seat_id"`
		SeatNumbers string `db:"seat_number"`
	}
	var S []SeatDetails
	q2 := `select s.id as seat_id,s.seat_number  from seat_bookings as sb join seats as s on sb.seat_id=s.id where sb.booking_id=$1 order by(sb.id,s.id);`
	err = r.db.SelectContext(c, &S, q2, bookingID)
	if err != nil {
		return nil, err
	}
	if len(S) == 0 {
		return nil, errors.New("no seat booked yet")
	}

	query := `SELECT 
    b.id AS booking_id,
    u.id AS user_id,
    u.username,
    sh.id AS show_id,
    m.title AS movie_name,
    t.theater_name,
    h.hall_name,
    b.status,
    b.total_amount,
    b.created_at
FROM bookings AS b
JOIN users AS u 
    ON b.user_id = u.id
JOIN shows AS sh 
    ON b.show_id = sh.id
JOIN movies AS m 
    ON m.id = sh.movie_id
JOIN halls AS h 
    ON h.id = sh.hall_id
JOIN theaters AS t 
    ON t.id = h.theater_id

WHERE b.id = $1
ORDER BY 
    t.id,
    h.id,
    sh.id;`

	var BookingDetails models.BookingsDetails
	err = r.db.GetContext(c, &BookingDetails, query, bookingID)
	if err != nil {
		return nil, err
	}

	for _, s := range S {
		BookingDetails.SeatID = append(BookingDetails.SeatID, s.SeatIds)
		BookingDetails.SeatNumber = append(BookingDetails.SeatNumber, s.SeatNumbers)

	}

	return &BookingDetails, nil

}

func (r *TheaterRepository) UserBookingHistoryRepo(c context.Context, userID int) (string, []models.BookingHistory, error) {

	//check if the user is exists or not
	query1 := `select username from users where id=$1`
	var UserName string
	err := r.db.GetContext(c, &UserName, query1, userID)
	if err != nil {
		return "", nil, err
	}
	if UserName == "" {
		return "", nil, errors.New("user not found")
	}

	query2 := `SELECT b.id AS booking_id,m.title AS movie_name,t.theater_name,h.hall_name,COUNT(sb.id) AS seat_count,b.status,b.total_amount,b.created_at FROM users AS u
               JOIN bookings AS b ON b.user_id = u.id
               JOIN seat_bookings AS sb ON sb.booking_id = b.id
               JOIN shows AS sh ON sh.id = b.show_id
               JOIN movies AS m ON m.id = sh.movie_id
               JOIN halls AS h ON h.id = sh.hall_id
               JOIN theaters AS t ON t.id = h.theater_id
               WHERE u.id = $1
               GROUP BY 
			            b.id,
			            m.title,
                        t.theater_name,
                        h.hall_name,
                        b.status,
                        b.total_amount,
                        b.created_at
              ORDER BY b.created_at DESC;`

	var BookingHstry []models.BookingHistory
	err = r.db.SelectContext(c, &BookingHstry, query2, userID)
	if err != nil {
		return "", nil, err
	}
	if len(BookingHstry) == 0 {
		return UserName, nil, errors.New("no booking history found")
	}

	return UserName, BookingHstry, nil

}

func (r *TheaterRepository) BookingCancelRepo(c context.Context, BookingID int) (*models.Bookings, error) {

	//check if the user is exists or not

	type BookingAndStatusCheck struct {
		Id     int    `db:"id"`
		Status string `db:"status"`
	}
	var BASC BookingAndStatusCheck

	query1 := `select id,status from bookings where id=$1`

	err := r.db.GetContext(c, &BASC, query1, BookingID)
	if err != nil {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("booking id not found")
	} else if BASC.Status == "cancelled" {
		return nil, errors.New("booking is already cancelled")
	}

	q2 := `update bookings 
	set 
	   status='cancelled',
	   updated_at=NOW()
    where id=$1  
	returning id as booking_id,
	          user_id,
              show_id,
              status,
              total_amount,
              created_at,
              updated_at;`

	var bookingDetails models.Bookings
	err = r.db.GetContext(c, &bookingDetails, q2, BookingID)

	if err != nil {
		return nil, err
	}
	return &bookingDetails, nil

}
