package models

type Theater struct {
	Id            int    `db:"id"`
	TheaterName   string `db:"theater_name"`
	NumberOfHalls int    `db:"number_of_halls"`
	City          string `db:"city"`
	State         string `db:"state"`
}
