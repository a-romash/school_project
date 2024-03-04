package model

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login          string `json:"login"`
	Name           string `json:"name"`
	Lastname       string `json:"lastname"`
	School         string `json:"school"`
	hashedPassword []byte
}

func CreateUser(login, name, lastname, school, password string) (u *User, err error) {
	hPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u = &User{
		Login:          login,
		Name:           name,
		Lastname:       lastname,
		School:         school,
		hashedPassword: hPassword,
	}
	return u, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.hashedPassword, []byte(password))
	// If password doesn't match, err will be non-nil and it returns false,
	// else password is correct and it return true
	return err == nil
}

func (u *User) GetJson() (data []byte, err error) {
	data, err = json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return data, nil
}
