package auth

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type AuthService interface {
	CreateUserService(c context.Context, useData dto.Users) error
	GetUserService(c context.Context, userID int) (dto.Users, error)
}
