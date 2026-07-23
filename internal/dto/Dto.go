package dto

import "time"

type Theater struct {
	Id            int    `json:"id"`
	TheaterName   string `json:"theater_name"`
	NumberOfHalls int    `json:"number_of_halls"`
	City          string `json:"city"`
	State         string `json:"state"`
}

type TheaterShows struct {
	ShowId      int       `json:"show_id"`
	TheaterName string    `json:"theater_name"`
	HallName    string    `json:"hall_name"`
	MovieName   string    `json:"movie_name"`
	Language    string    `json:"language"`
	StartDate   time.Time `json:"starts_at"`
	EndDate     time.Time `json:"ends_at"`
	Price       float64   `json:"base_price"`
}

type SeatsInShows struct {
	ShowId         int     `json:"show_id"`
	MovieName      string  `json:"movie_name"`
	HallName       string  `json:"hall_name"`
	SeatsAvailable []Seats `json:"seats"`
}

type Seats struct {
	SeatId     int    `json:"seat_id"`
	SeatNumber string `json:"seat_number"`
	SeatType   string `json:"seat_type"`
	Status     string `json:"status"`
}

// type SeatBooking struct {
// 	Id int `json:"id"`
// 	SeatID int `json:"seat_id"`
// 	ShowID int `json:"show_id"`
// }

// create table if not exists seat_bookings(
// id bigserial primary key,
// booking_id bigint not null references bookings(id),
// seat_id bigint not null references seats(id),
// show_id bigint not null references shows(id),
// created_at  timestamptz not null default now(),
// unique(show_id,seat_id)
// );
