package postgresql

import (
	"context"
	"errors"
	"project/pkg/model"
	"time"

	"github.com/google/uuid"
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

func (db *Postgresql) NewToken(login string) (token *model.Token, err error) {
	const sql = `
	INSERT INTO tokens (login, token, expires_at)
  	VALUES ($1, $2, $3);
	`

	token = model.CreateToken(login)

	_, err = db.pool.Exec(context.Background(), sql, token.Login, token.Token, token.Expires_at)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (db *Postgresql) DeleteToken(token uuid.UUID) error {
	const sql = `
	DELETE FROM tokens
	WHERE token = $1;
	`

	_, err := db.pool.Exec(context.Background(), sql, token)
	if err != nil {
		return err
	}
	return nil
}

func (db *Postgresql) DeleteTokensByLogin(login string) error {
	const sql = `
	DELETE FROM tokens
	WHERE login = $1;
	`

	_, err := db.pool.Exec(context.Background(), sql, login)
	if err != nil {
		return err
	}
	return nil
}

type ErrTokenExpired error

func (db *Postgresql) ValidateToken(token uuid.UUID) error {
	const sql = `
	SELECT *
	FROM tokens
	WHERE token = $1;
	`

	rows, _ := db.pool.Query(context.Background(), sql, token)
	tokens, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Token])
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return pgx.ErrNoRows
	}
	t := &tokens[0]
	if time.Now().After(t.Expires_at) {
		db.DeleteToken(token)
		var e ErrTokenExpired = errors.New("token is expired")
		return e
	}
	return nil
}
