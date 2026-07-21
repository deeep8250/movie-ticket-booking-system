package models

import "time"

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
