package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/deeep8250/movie-ticket-booking-system/internal/config"
	"github.com/deeep8250/movie-ticket-booking-system/internal/dto"
	"github.com/deeep8250/movie-ticket-booking-system/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo() *AuthRepo {
	return &AuthRepo{
		db: config.DBClients.PostgresClient,
	}
}

func (r *AuthRepo) CreateUserRepo(c context.Context, userData dto.UsersRequest) error {

	query := `insert into users(username,email,mobile,password_hash) values($1,$2,$3,$4)`
	result, err := r.db.ExecContext(c, query, userData.UserName, userData.Email, userData.Mobile, userData.Password)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("unable to create user")
	}

	return nil

}

// clean ups
func (r *AuthRepo) GetUserRepo(c context.Context, userID int) (*models.Users, error) {

	query := `select id,username,email,mobile,created_at,updated_at from users where id=$1`

	var user models.Users
	err := r.db.GetContext(c, &user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil

}

func (r *AuthRepo) GetUserByEmailRepo(c context.Context, Email string) (*models.Users, error) {
	query := `select * from users where email=$1`

	var user models.Users
	err := r.db.GetContext(c, &user, query, Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}
	return &user, nil

}

func (r *AuthRepo) VerifyEmail(c context.Context, Email string) (bool, error) {
	query := `select count(*) from users where email=$1`

	var userCount int
	err := r.db.GetContext(c, &userCount, query, Email)
	if err != nil {

		return false, err
	}
	return userCount > 0, nil
}
func (r *AuthRepo) VerifyMobile(c context.Context, mobile string) (bool, error) {
	query := `select count(*) from users where mobile=$1`

	var userCount int
	err := r.db.GetContext(c, &userCount, query, mobile)
	if err != nil {

		return false, err
	}

	return userCount > 0, nil

}
