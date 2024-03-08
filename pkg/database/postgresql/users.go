package postgresql

import (
	"context"
	"project/pkg/model"

	"github.com/jackc/pgx/v5"
)

func (db *Postgresql) RegisterNewUser(user *model.User) error {
	const sql = `
	INSERT INTO users (login, name, lastname, school, hashedPassword)
  	VALUES ($1, $2, $3, $4, $5);
	`

	_, err := db.pool.Exec(context.Background(), sql, user.Login, user.Name, user.Lastname, user.School, user.HashedPassword)
	return err
}

func (db *Postgresql) GetUser(login string) (user *model.User, err error) {
	const sql = `
	SELECT * FROM users WHERE login = $1
	`

	rows, _ := db.pool.Query(context.Background(), sql, login)
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, pgx.ErrNoRows
	}
	user = &users[0]
	return user, nil
}
