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

type SeatBooking struct {
	Id          int     `json:"booking_id"`
	UserId      int     `json:"user_id"`
	ShowID      int     `json:"show_id"`
	SeatsBooked []int   `json:"seats_ids"`
	TotalPrice  float64 `json:"total_price"`
}

type BookingsDetails struct {
	BookingID   int       `json:"booking_id"`
	UserID      int       `json:"user_id"`
	UserName    string    `json:"username"`
	ShowID      int       `json:"show_id"`
	MovieName   string    `json:"movie_name"`
	TheaterName string    `json:"theater_name"`
	HallName    string    `json:"hall_name"`
	SeatID      []int     `json:"seat_id"`
	SeatNumber  []string  `json:"seat_number"`
	Status      string    `json:"status"`
	TotalAmount int       `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type BookingHistory struct {
	BookingID   int       `json:"booking_id"`
	MovieName   string    `json:"movie_name"`
	TheaterName string    `json:"theater_name"`
	HallName    string    `json:"hall_name"`
	SeatCount   int       `json:"seat_count"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type BookingHistoryWithUsrDtlsAtchd struct {
	UserID          int              `json:"user_id"`
	Username        string           `json:"username"`
	BookingHitories []BookingHistory `json:"bookings"`
}

type BookingsDto struct {
	BookingID   int       `json:"booking_id"`
	UserID      int       `json:"user_id"`
	ShowID      int       `json:"show_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TotalAmount float64   `json:"total_amount"`
}

type GetMovies struct {
	MoviesID    int       `json:"movie_id"`
	Title       string    `json:"title"`
	Language    string    `json:"language"`
	DurationMin int       `json:"duration_min"`
	ReleaseDate time.Time `json:"release_date"`
}

type Shows struct {
	ShowID      int       `json:"show_id"`
	TheaterName string    `json:"theater_name"`
	HallName    string    `json:"halll_name"`
	City        string    `json:"city"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	BasePrice   float64   `json:"base_price"`
}

type UsersRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email" binding:"required,email"`
	Mobile   int    `json:"mobile" binding:"required,min=10"`
	Password string `json:"password" binding:"required,min=5"`
}

type Users struct {
	UserID       int       `json:"id"`
	UserName     string    `json:"username"`
	Email        string    `json:"email" `
	Mobile       int       `json:"mobile" `
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	PasswordHash string    `json:"-" `
}
