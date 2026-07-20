package dto

type Theater struct {
	Id            int    `json:"id"`
	TheaterName   string `json:"theater_name"`
	NumberOfHalls int    `json:"number_of_halls"`
	City          string `json:"city"`
	State         string `json:"state"`
}
