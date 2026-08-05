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

func (r *AuthRepo) CreateUserRepo(c context.Context, userData dto.Users) error {

	query := `insert into users(username,email,mobile,password_hash) values($1,$2,$3,$4)`
	result, err := r.db.ExecContext(c, query, userData.UserName, userData.Email, userData.Mobile, userData.PasswordHash)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("unable to crete user")
	}

	return nil

}
func (r *AuthRepo) GetUserRepo(c context.Context, userID int) (*models.Users, error) {

	query := `select * from users where id=$1`

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

func (r *AuthRepo) GetUserByIdRepo(c context.Context, userID int) (string, error) {
	q := `select * from users where id=$1`
	var user models.Users
	err := r.db.GetContext(c, &user, q, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	return user.PasswordHash, nil
}
