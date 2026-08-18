package shows

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type ShowServices struct {
	repo ShowsRepositoriesInterface
}

// theaters handlers
func NewShowsServices(r ShowsRepositoriesInterface) *ShowServices {
	return &ShowServices{
		repo: r,
	}
}

func (s *ShowServices) GetShowsService(c context.Context, TheaterId int) ([]dto.TheaterShows, error) {
	shows, err := s.repo.GetShowsRepo(c, TheaterId)
	if err != nil {
		return nil, err
	}

	var showsInTheater []dto.TheaterShows
	for r := range shows {

		showInTheater := dto.TheaterShows{
			ShowId:      shows[r].ShowId,
			TheaterName: shows[r].TheaterName,
			HallName:    shows[r].HallName,
			MovieName:   shows[r].MovieName,
			Language:    shows[r].Language,
			StartDate:   shows[r].StartDate,
			EndDate:     shows[r].EndDate,
			Price:       shows[r].Price,
		}
		showsInTheater = append(showsInTheater, showInTheater)

	}

	return showsInTheater, nil
}

func (s *ShowServices) GetShowsByMovieIDService(c context.Context, MovieID int) ([]dto.Shows, error) {

	sh, err := s.repo.GetShowsByMovieIDRepo(c, MovieID)
	if err != nil {
		return nil, err
	}

	var shows []dto.Shows
	for _, show := range sh {
		s := dto.Shows{
			ShowID:      show.ShowID,
			TheaterName: show.TheaterName,
			City:        show.City,
			HallName:    show.HallName,
			StartsAt:    show.StartsAt,
			EndsAt:      show.EndsAt,
			BasePrice:   show.BasePrice,
		}
		shows = append(shows, s)
	}

	return shows, nil

}
