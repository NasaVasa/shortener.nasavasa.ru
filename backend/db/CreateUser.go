package db

func (pgContext *PgContext) CreateUser(login string, password string, email string, token string) error {
	query, args, err := pgContext.Psql.Insert("users").
		Columns("login", "password", "email", "token").
		Values(login, password, email, token).
		ToSql()
	if err != nil {
		return err
	}
	_, err = pgContext.DB.Exec(query, args...)
	return err
}
