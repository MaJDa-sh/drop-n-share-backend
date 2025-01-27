package services

import (
	"drop-n-share/internal/models"
	"errors"
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/jmoiron/sqlx"
)

type UserService struct {
	db *sqlx.DB
}

func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{db: db}
}

func (u *UserService) SignUp(user *models.SignParams) (*models.User, error) {
	var result models.User

	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	err = u.db.Get(&result, `insert into users (username, password_hash) values ($1, $2) returning *`, user.Username, hash)
	if err != nil {
		return nil, errors.New("user already exists")
	}
	return &result, nil
}

func (u *UserService) SignIn(user *models.SignParams) (*models.User, error) {
	var result models.User
	err := u.db.Get(&result, `select * from users where username = $1`, user.Username)
	if err != nil {
		return nil, err
	}

	if _, err := result.CheckPassword(user.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &result, nil
}

func (us *UserService) GetUserByID(userID int) (*models.UserWithFiles, error) {
	var user models.UserWithFiles

	userQuery := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id = $1`
	err := us.db.Get(&user, userQuery, userID)
	if err != nil {
		log.Println("Error fetching user:", err)
		return nil, errors.New("user not found")
	}

	filesQuery := `SELECT id, file_name, category, user_id, created_at 
				   FROM files WHERE user_id = $1`
	var files []models.File
	err = us.db.Select(&files, filesQuery, userID)
	if err != nil {
		log.Println("Error fetching files:", err)

	}

	user.Files = files

	return &user, nil
}
