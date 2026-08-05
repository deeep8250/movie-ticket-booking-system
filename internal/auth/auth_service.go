package auth

import (
	"context"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo AuthRepoInterface
}

func NewAuthService(Repo AuthRepoInterface) *AuthService {
	return &AuthService{
		repo: Repo,
	}
}

func (s *AuthService) CreateUserService(c context.Context, useData dto.Users) error {

	PassHash, err := bcrypt.GenerateFromPassword([]byte(useData.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	useData.PasswordHash = string(PassHash)

	err = s.repo.CreateUserRepo(c, useData)
	if err != nil {
		return err
	}

	return nil

}

func (s *AuthService) GetUserService(c context.Context, userID int) (*dto.Users, error) {
	user, err := s.repo.GetUserRepo(c, userID)
	if err != nil {
		return nil, err
	}
	UserDTO := dto.Users{
		UserID:       user.UserID,
		UserName:     user.UserName,
		Email:        user.Email,
		Mobile:       user.Mobile,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		PasswordHash: user.PasswordHash,
	}
	return &UserDTO, nil
}
