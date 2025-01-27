package views

import "drop-n-share/internal/models"

type FilesView struct {
	Files []models.File `json:"files"`
}
