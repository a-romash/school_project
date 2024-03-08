package model

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Login      string    `db:"login"`
	Token      uuid.UUID `db:"t"`
	Expires_at time.Time `db:"expires_at"`
}

// expires_after = 43200 secs (12 hours)
func CreateToken(login string) (token *Token) {
	token = CreateTokenWithExpire(login, 43200)
	return token
}

// expires_after in seconds
func CreateTokenWithExpire(login string, expires_after int) (token *Token) {
	token = &Token{
		Login:      login,
		Token:      uuid.New(),
		Expires_at: time.Now().Add(time.Second * time.Duration(expires_after)),
	}
	return token
}
