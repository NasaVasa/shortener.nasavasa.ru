package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
)

type User struct {
	Id             pgtype.UUID
	Login          string
	Email          string
	Password       string
	Token          string
	TokenCreatedAt pgtype.Timestamp
	UserCreatedAt  pgtype.Timestamp
}

// NotFoundUserError is error that returns when user not found
type NotFoundUserError struct{}

func (e *NotFoundUserError) Error() string {
	return "user not found"
}

// GetUser returns user by login using squirrel
func (pgContext *PgContext) GetUser(getBy string, value string) (User, error) {
	var user User
	query, args, err := pgContext.Psql.Select("id", "login", "email", "password", "token", "token_created_at", "user_created_at").
		From("users").
		Where(sq.Eq{getBy: value}).ToSql()
	if err != nil {
		return user, err
	}
	err = pgContext.DB.QueryRow(query, args...).Scan(&user.Id, &user.Login, &user.Email, &user.Password, &user.Token, &user.TokenCreatedAt, &user.UserCreatedAt)
	if err != nil {
		return user, &NotFoundUserError{}
	}
	return user, nil
}
