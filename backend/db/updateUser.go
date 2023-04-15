package db

import sq "github.com/Masterminds/squirrel"

func (pgContext *PgContext) UpdateUser(user User) error {
	query, args, err := pgContext.Psql.Update("users").
		Set("token", user.Token).
		Set("token_created_at", user.TokenCreatedAt).
		Where(sq.Eq{"login": user.Login}).ToSql()
	if err != nil {
		return err
	}
	_, err = pgContext.DB.Exec(query, args...)
	return err
}
