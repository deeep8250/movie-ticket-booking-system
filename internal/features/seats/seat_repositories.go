package seats

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SeatsRepositories struct {
	db    *sqlx.DB
	redis *redis.Client
}

// theaters handlers
func NewSeatsRepositories() *SeatsRepositories {
	return &SeatsRepositories{
		db:    config.DBClients.PostgresClient,
		redis: config.DBClients.RedisClient,
	}
}

// seats
func (r *SeatsRepositories) GetSeats(c context.Context, showsId int) (*models.SeatsInShows, error) {

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

	for i, seat := range Seats {
		if seat.Status == "booked" || seat.Status == "disable" {
			continue
		}
		key := fmt.Sprintf("seat_lock:show:%d:seat:%d", showsId, seat.SeatId)

		exists, err := r.redis.Exists(c, key).Result()
		if err != nil {
			return nil, err
		}
		if exists > 0 {
			Seats[i].Status = "locked"
		}
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
func (r *SeatsRepositories) SeatLockRepo(c context.Context, userID int, showID int, seatIDs []int) error {

	query := `select count(*) from seats as s join shows as sh on sh.hall_id=s.hall_id where s.id=ANY($1) and sh.id=$2`
	var seatCount int
	err := r.db.GetContext(c, &seatCount, query, pq.Array(seatIDs), showID)
	if err != nil {
		return err
	}
	if seatCount != len(seatIDs) {
		return errors.New("one or more seat ids do not belong to this show")
	}

	// checking if the selected seats are already booked or not
	query2 := `SELECT count(*)    
               FROM seat_bookings AS sb
               JOIN bookings AS b
               ON b.id = sb.booking_id
               WHERE b.show_id = $1
               AND sb.seat_id = ANY($2)
               AND b.status = 'confirmed';`

	// i want to know if the selected seats are booked or not
	var count int
	err = r.db.GetContext(c, &count, query2, showID, pq.Array(seatIDs))
	if err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("seat already booked by other user %d", count)
	}
	for _, seatID := range seatIDs {
		key := fmt.Sprintf("seat_lock:show:%d:seat:%d", showID, seatID)

		ok, err := r.redis.SetNX(c, key, strconv.Itoa(userID), 5*time.Minute).Result()
		if err != nil {
			return err
		}

		if !ok {
			lockOwner, err := r.redis.Get(c, key).Result()
			if err != nil {
				return err
			}

			if lockOwner == strconv.Itoa(userID) {
				return errors.New("seat already locked by you")
			}
			return errors.New("seat already locked by another user")
		}
	}
	return nil
}
func (r *SeatsRepositories) SeatUnLockRepo(c context.Context, userID int, showID int, seatIDs []int) error {

	query := `select count(*) from seats as s join shows as sh on sh.hall_id=s.hall_id where s.id=ANY($1) and sh.id=$2`
	var seatCount int
	err := r.db.GetContext(c, &seatCount, query, pq.Array(seatIDs), showID)
	if err != nil {
		return fmt.Errorf("failed to validate seats for show: %w", err)
	}
	if seatCount != len(seatIDs) {
		return errors.New("one or more seat ids do not belong to this show")
	}

	//

	for _, seatID := range seatIDs {
		key := fmt.Sprintf("seat_lock:show:%d:seat:%d", showID, seatID)

		result, err := r.redis.Get(c, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return errors.New("seat is already unlocked")
			}
			return err
		}
		if result != strconv.Itoa(userID) {
			return errors.New("seat is locked by another user")
		}

		if err := r.redis.Del(c, key).Err(); err != nil {
			return err
		}

	}
	return nil
}
