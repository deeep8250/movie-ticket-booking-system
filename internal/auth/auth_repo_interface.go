package auth

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type AuthRepoInterface interface {
	CreateUserRepo(c context.Context, useData dto.Users) error
	GetUserRepo(c context.Context, userID int) (models.Users, error)
	GetUserByIdRepo(c context.Context, userID int) (string, error)
}
