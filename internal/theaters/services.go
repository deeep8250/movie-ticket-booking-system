package theaters

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type TheaterService struct {
	repo TheaterRepoInterface
}

func NewTheaterService(Repo *TheaterRepository) *TheaterService {
	return &TheaterService{
		repo: Repo,
	}
}

func (s *TheaterService) GetTheatersService(c context.Context) ([]dto.Theater, error) {
	theatersValues, err := s.repo.GetTheaters(c)
	if err != nil {
		return nil, err
	}

	var theaters []dto.Theater
	for r := range theatersValues {

		theater := dto.Theater{
			Id:            theatersValues[r].Id,
			TheaterName:   theatersValues[r].TheaterName,
			NumberOfHalls: theatersValues[r].NumberOfHalls,
			City:          theatersValues[r].City,
			State:         theatersValues[r].State,
		}
		theaters = append(theaters, theater)

	}

	return theaters, nil
}

func (s *TheaterService) GetShowsService(c context.Context, TheaterId int) ([]dto.TheaterShows, error) {
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
