package services

import (
	"drop-n-share/internal/models"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type FileService struct {
	db *sqlx.DB
}

func NewFileService(db *sqlx.DB) *FileService {
	return &FileService{
		db: db,
	}
}

func (fs *FileService) SaveFile(file models.File) (models.File, error) {
	query := `INSERT INTO files (file_name, category, user_id, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`

	var fileID int
	err := fs.db.QueryRow(query, file.FileName, file.Category, file.UserID, file.CreatedAt).Scan(&fileID)
	if err != nil {
		log.Println("Error saving file:", err)
		return models.File{}, err
	}
	file.ID = fileID
	return file, nil
}

func (fs *FileService) GetFileByID(file *models.File, fileID string, userID uint) error {

	// query := `SELECT id, file_name FROM files WHERE id = $1 and user_id = $2`

	query := `SELECT id, file_name FROM files WHERE id = $1`
	err := fs.db.Get(file, query, fileID)
	if err != nil {
		log.Println("Error retrieving file by ID:", err)
		return err
	}
	return nil
}

func (fs *FileService) AllFiles() ([]models.File, error) {
	query := `SELECT id, file_name, category, user_id FROM files`
	var files []models.File

	err := fs.db.Select(&files, query)
	if err != nil {
		log.Println("Error retrieving files:", err)
		return nil, err
	}

	return files, nil
}

func (fs *FileService) FindFiles(ids []uint) ([]models.File, error) {
	query := `SELECT id, file_name, category, user_id FROM files WHERE id = ANY($1)`
	var files []models.File

	err := fs.db.Select(&files, query, pq.Array(ids))
	if err != nil {
		log.Println("Error retrieving files:", err)
		return nil, err
	}

	return files, nil
}
