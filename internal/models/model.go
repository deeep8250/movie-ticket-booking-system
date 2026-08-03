package models

import (
	"time"
)

type Theater struct {
	Id            int    `db:"id"`
	TheaterName   string `db:"theater_name"`
	NumberOfHalls int    `db:"number_of_halls"`
	City          string `db:"city"`
	State         string `db:"state"`
}

type TheaterShows struct {
	ShowId      int       `db:"show_id"`
	TheaterName string    `db:"theater_name"`
	HallName    string    `db:"hall_name"`
	MovieName   string    `db:"movie_name"`
	Language    string    `db:"language"`
	StartDate   time.Time `db:"starts_at"`
	EndDate     time.Time `db:"ends_at"`
	Price       float64   `db:"base_price"`
}

type SeatsInShows struct {
	ShowId         int     `db:"show_id"`
	MovieName      string  `db:"movie_name"`
	HallName       string  `db:"hall_name"`
	SeatsAvailable []Seats `db:"seats"`
}

type Seats struct {
	SeatId     int    `db:"seat_id"`
	SeatNumber string `db:"seat_number"`
	SeatType   string `db:"seat_type"`
	Status     string `db:"status"`
}

type SeatBooking struct {
	Id          int     `db:"id"`
	UserId      int     `db:"user_id"`
	ShowID      int     `db:"show_id"`
	SeatsBooked []int   `db:"seats_booked"`
	TotalPrice  float64 `db:"total_amount"`
}

type BookingsDetails struct {
	BookingID   int       `db:"booking_id"`
	UserID      int       `db:"user_id"`
	UserName    string    `db:"username"`
	ShowID      int       `db:"show_id"`
	MovieName   string    `db:"movie_name"`
	TheaterName string    `db:"theater_name"`
	HallName    string    `db:"hall_name"`
	SeatID      []int     `db:"seat_id"`
	SeatNumber  []string  `db:"seat_number"`
	Status      string    `db:"status"`
	TotalAmount float64   `db:"total_amount"`
	CreatedAt   time.Time `db:"created_at"`
}

type BookingHistory struct {
	BookingID   int       `db:"booking_id"`
	Username    string    `db:"username"`
	MovieName   string    `db:"movie_name"`
	TheaterName string    `db:"theater_name"`
	HallName    string    `db:"hall_name"`
	SeatCount   int       `db:"seat_count"`
	Status      string    `db:"status"`
	TotalAmount float64   `db:"total_amount"`
	CreatedAt   time.Time `db:"created_at"`
}

type Bookings struct {
	BookingID   int       `db:"booking_id"`
	UserID      int       `db:"user_id"`
	ShowID      int       `db:"show_id"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	TotalAmount float64   `db:"total_amount"`
}

type GetMovies struct {
	MoviesID    int       `db:"id"`
	Title       string    `db:"title"`
	Language    string    `db:"language"`
	DurationMin int       `db:"duration_min"`
	ReleaseDate time.Time `db:"release_date"`
}
