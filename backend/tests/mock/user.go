package mock

import (
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"shortener/db"
	"time"
)

type User struct {
	id             uuid.UUID
	login          string
	email          string
	password       string
	Token          string
	TokenCreatedAt time.Time
	userCreatedAt  time.Time
}

func NewUser(id uuid.UUID, login string, email string, password string, token string, tokenCreatedAt time.Time, userCreatedAt time.Time) User {
	if id == uuid.Nil {
		id = uuid.New()
	}
	if tokenCreatedAt.IsZero() {
		tokenCreatedAt = time.Now()
	}
	if userCreatedAt.IsZero() {
		userCreatedAt = time.Now()
	}
	return User{
		id:             id,
		login:          login,
		email:          email,
		password:       password,
		Token:          token,
		TokenCreatedAt: tokenCreatedAt,
		userCreatedAt:  userCreatedAt,
	}
}

func (pgMock *PgMock) CheckToken(token string) bool {
	for _, user := range pgMock.Users {
		if user.Token == token {
			return true
		}
	}
	return false
}

func (pgMock *PgMock) CheckUser(getBy string, value string) bool {
	for _, user := range pgMock.Users {
		if getBy == "Id" {
			if user.id.String() == value {
				return true
			}
		}
		if getBy == "login" {
			if user.login == value {
				return true
			}
		}
		if getBy == "email" {
			if user.email == value {
				return true
			}
		}
	}
	return false
}

func (pgMock *PgMock) CreateUser(login string, password string, email string, token string) error {
	pgMock.Users = append(pgMock.Users, NewUser(uuid.New(), login, email, password, token, time.Now(), time.Now()))
	return nil
}

func (pgMock *PgMock) GetUser(getBy string, value string) (db.User, error) {
	for _, user := range pgMock.Users {
		if getBy == "Id" {
			if user.id.String() == value {
				var id [16]byte
				copy(id[:], user.id[:])
				// parse user to db.User
				return db.User{
					Id:             pgtype.UUID{Bytes: id, Status: pgtype.Present},
					Login:          user.login,
					Email:          user.email,
					Password:       user.password,
					Token:          user.Token,
					TokenCreatedAt: pgtype.Timestamp{Time: user.TokenCreatedAt, Status: pgtype.Present},
					UserCreatedAt:  pgtype.Timestamp{Time: user.userCreatedAt, Status: pgtype.Present},
				}, nil
			}
		}
		if getBy == "login" {
			if user.login == value {
				var id [16]byte
				copy(id[:], user.id[:])
				// parse user to db.User
				return db.User{
					Id:             pgtype.UUID{Bytes: id, Status: pgtype.Present},
					Login:          user.login,
					Email:          user.email,
					Password:       user.password,
					Token:          user.Token,
					TokenCreatedAt: pgtype.Timestamp{Time: user.TokenCreatedAt, Status: pgtype.Present},
					UserCreatedAt:  pgtype.Timestamp{Time: user.userCreatedAt, Status: pgtype.Present},
				}, nil
			}
		}
		if getBy == "email" {
			if user.email == value {
				var id [16]byte
				copy(id[:], user.id[:])
				// parse user to db.User
				return db.User{
					Id:             pgtype.UUID{Bytes: id, Status: pgtype.Present},
					Login:          user.login,
					Email:          user.email,
					Password:       user.password,
					Token:          user.Token,
					TokenCreatedAt: pgtype.Timestamp{Time: user.TokenCreatedAt, Status: pgtype.Present},
					UserCreatedAt:  pgtype.Timestamp{Time: user.userCreatedAt, Status: pgtype.Present},
				}, nil
			}
		}
		if getBy == "token" {
			if user.Token == value {
				var id [16]byte
				copy(id[:], user.id[:])
				// parse user to db.User
				return db.User{
					Id:             pgtype.UUID{Bytes: id, Status: pgtype.Present},
					Login:          user.login,
					Email:          user.email,
					Password:       user.password,
					Token:          user.Token,
					TokenCreatedAt: pgtype.Timestamp{Time: user.TokenCreatedAt, Status: pgtype.Present},
					UserCreatedAt:  pgtype.Timestamp{Time: user.userCreatedAt, Status: pgtype.Present},
				}, nil
			}
		}
	}
	return db.User{}, &db.NotFoundUserError{}
}

func (pgMock *PgMock) UpdateUser(user db.User) error {
	for i, u := range pgMock.Users {
		if u.login == user.Login {
			pgMock.Users[i].Token = user.Token
			pgMock.Users[i].TokenCreatedAt = user.TokenCreatedAt.Time
		}
	}
	return nil
}

func (pgMock *PgMock) UpdateUserTokenTime(token string) error {
	for i, u := range pgMock.Users {
		if u.Token == token {
			pgMock.Users[i].TokenCreatedAt = time.Now()
		}
	}
	return nil
}
