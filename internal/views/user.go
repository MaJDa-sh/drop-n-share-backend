package views

import "drop-n-share/internal/models"

type NewUserData struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Token    string `json:token`
}

type ProfileData struct {
	ID       uint          `json:"id"`
	Username string        `json:"username"`
	Files    []models.File `json:"files"`
}

func NewUserView(id uint, username string, token string) *NewUserData {
	return &NewUserData{
		ID:       id,
		Username: username,
		Token:    token,
	}
}

func ProfileView(id uint, username string, files []models.File) *ProfileData {
	return &ProfileData{
		ID:       id,
		Username: username,
		Files:    files,
	}
}
