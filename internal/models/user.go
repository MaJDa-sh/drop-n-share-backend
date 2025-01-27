// internal/models/user.go
package models

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
)

type User struct {
	ID           uint      `db:"id"`
	Username     string    `db:"username"`
	PasswordHash *string   `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserWithFiles struct {
	ID           uint      `db:"id"`
	Username     string    `db:"username"`
	PasswordHash *string   `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Files        []File    `db:"files"`
}
type SignParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var ErrUserAlreadyExists = errors.New("user with this username already exists")
var ErrInvalidPassword = errors.New("password is invalid")
var ErrPasswordsDontMatch = errors.New("passwords don't match")

func (u *User) CheckPassword(password string) (bool, error) {
	if !strings.HasPrefix(*u.PasswordHash, "$argon2id$") {
		return false, ErrInvalidPassword
	}

	match, err := argon2id.ComparePasswordAndHash(password, *u.PasswordHash)

	if err != nil {
		log.Println("Check password error ", err)
		return false, nil
	}

	if !match {
		return false, ErrPasswordsDontMatch
	}

	return true, nil
}
