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
