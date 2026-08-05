package auth

import (
	"context"
	"errors"

	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
	"github.com/deeep8250/movie-ticket-booking-system/internal/jwtToken"
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

func (s *AuthService) CreateUserService(c context.Context, useData dto.UsersRequest) error {

	isEmailExists, err := s.repo.VerifyEmail(c, useData.Email)
	if err != nil {
		return err
	}
	if isEmailExists {
		return errors.New("email already exists")
	}

	isMobExists, err := s.repo.VerifyMobile(c, useData.Mobile)
	if err != nil {
		return err
	}
	if isMobExists {
		return errors.New("mobile already exists")
	}

	PassHash, err := bcrypt.GenerateFromPassword([]byte(useData.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	useData.Password = string(PassHash)

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
		UserID:    user.UserID,
		UserName:  user.UserName,
		Email:     user.Email,
		Mobile:    user.Mobile,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return &UserDTO, nil
}

func (s *AuthService) LoginService(c context.Context, cred dto.UserLogin) (string, error) {

	user, err := s.repo.GetUserByEmailRepo(c, cred.Email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(cred.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := jwtToken.CreateJWT(user.UserID)
	if err != nil {
		return "", err
	}
	return token, nil

}
