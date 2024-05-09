package repository

import (
	"context"
	"database/sql"
	"projectsphere/eniqlo-store/internal/staff/entity"
	"projectsphere/eniqlo-store/pkg/database"
	"projectsphere/eniqlo-store/pkg/protocol/msg"
	"strings"
)

type UserRepo struct {
	dbConnector database.PostgresConnector
}

func NewUserRepo(dbConnector database.PostgresConnector) UserRepo {
	return UserRepo{
		dbConnector: dbConnector,
	}
}

func (r UserRepo) CreateUser(ctx context.Context, param entity.UserParam) (entity.User, error) {
	query := `
		INSERT INTO users (email, name, phone_number, password, salt) VALUES 
		($1, $2, $3, $4, $5) RETURNING user_id, email, name, phone_number, password, salt, created_at, updated_at
	`
	var row entity.User
	err := r.dbConnector.DB.GetContext(
		ctx,
		&row,
		query,
		param.Email,
		param.Name,
		param.PhoneNumber,
		param.Password,
		param.Salt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return entity.User{}, &msg.RespError{
				Code:    409,
				Message: msg.ErrEmailAlreadyExist,
			}
		} else {
			return entity.User{}, msg.InternalServerError(err.Error())
		}
	}

	return row, nil
}

func (r UserRepo) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (entity.User, error) {
	query := `
		SELECT user_id, email, name, phone_number, password, salt, created_at, updated_at FROM users WHERE phone_number = $1
	`

	var row entity.User
	err := r.dbConnector.DB.GetContext(
		ctx,
		&row,
		query,
		phoneNumber,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, msg.NotFound(msg.ErrUserNotFound)
		} else {
			return entity.User{}, msg.InternalServerError(err.Error())
		}
	}

	return row, nil
}

func (r UserRepo) IsUserExist(ctx context.Context, userId uint32) bool {
	query := `
		SELECT 1 FROM users WHERE id_user = $1
	`

	var result = 0
	err := r.dbConnector.DB.GetContext(
		ctx,
		&result,
		query,
		userId,
	)
	if err != nil {
		return false
	}

	return result == 1
}

func (r UserRepo) IsPhoneNumberExist(ctx context.Context, phoneNumber string) bool {
	query := `
		SELECT 1 FROM users WHERE phone_number = $1
	`

	var result = 0
	err := r.dbConnector.DB.GetContext(
		ctx,
		&result,
		query,
		phoneNumber,
	)

	if err != nil {
		return false
	}

	return result == 1
}
