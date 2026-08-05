package auth

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
)

type AuthServiceInterface interface {
	CreateUserService(c context.Context, useData dto.UsersRequest) error
	GetUserService(c context.Context, userID int) (dto.Users, error)
}
