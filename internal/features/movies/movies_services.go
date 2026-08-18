package movies

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type moviesServices struct {
	repo MoviesRepositoriesInterface
}

// theaters handlers
func NewMoviesServices(R MoviesRepositoriesInterface) *moviesServices {
	return &moviesServices{
		repo: R,
	}
}

func (s *moviesServices) GetMoviesService(c context.Context) ([]dto.GetMovies, error) {
	movies, err := s.repo.GetMoviesRepo(c)
	if err != nil {
		return nil, err
	}

	var ReturnedMovies []dto.GetMovies
	for _, movie := range movies {
		m := dto.GetMovies{
			MoviesID:    movie.MoviesID,
			Title:       movie.Title,
			Language:    movie.Language,
			DurationMin: movie.DurationMin,
			ReleaseDate: movie.ReleaseDate,
		}

		ReturnedMovies = append(ReturnedMovies, m)
	}

	return ReturnedMovies, nil

}
func (s *moviesServices) GetMovieByIDService(c context.Context, id int) (*dto.GetMovies, error) {

	m, err := s.repo.GetMovieByIDRepo(c, id)
	if err != nil {
		return nil, err
	}
	movie := dto.GetMovies{
		MoviesID:    m.MoviesID,
		Title:       m.Title,
		Language:    m.Language,
		DurationMin: m.DurationMin,
		ReleaseDate: m.ReleaseDate,
	}

	return &movie, nil
}
