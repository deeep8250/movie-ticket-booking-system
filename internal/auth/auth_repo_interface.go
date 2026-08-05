package auth

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
)

type AuthRepoInterface interface {
	CreateUserRepo(c context.Context, userData dto.UsersRequest) error
	GetUserRepo(c context.Context, userID int) (*models.Users, error)

	GetUserByEmailRepo(c context.Context, Email string) (*models.Users, error)
	VerifyEmail(c context.Context, Email string) (bool, error)
	VerifyMobile(c context.Context, mobile string) (bool, error)
}
