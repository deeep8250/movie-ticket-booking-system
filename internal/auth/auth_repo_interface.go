package auth

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type AuthRepo interface {
	CreateUserRepo(c context.Context, useData dto.Users) error
	GetUserRepo(c context.Context, userID int) (models.Users, error)
}
