package controllers

import (
	"drop-n-share/internal/middleware"
	"drop-n-share/internal/models"
	"drop-n-share/internal/services"
	"drop-n-share/internal/views"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FileController struct {
	FileService *services.FileService
}

func NewFileController(fileService *services.FileService) *FileController {
	return &FileController{
		FileService: fileService,
	}
}

func (fc *FileController) HandleFileUpload(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	category := r.FormValue("category")
	claims, err := middleware.GetClaimsFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid token"})
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid user ID in token"})
		return
	}

	if userID <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	fileName := fileHeader.Filename
	filePath := filepath.Join("uploads", fileName)

	outFile, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Saving file failed"})
		return
	}
	defer outFile.Close()

	_, err = outFile.ReadFrom(file)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to save the file"})
		return
	}

	newFile := models.File{
		FileName:  fileName,
		Category:  category,
		UserID:    uint(userID),
		CreatedAt: time.Now(),
	}

	savedFile, err := fc.FileService.SaveFile(newFile)
	if err != nil {
		http.Error(w, "Failed to save file metadata", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedFile)
}

func (fc *FileController) HandleFileDownload(w http.ResponseWriter, r *http.Request) {

	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	claims, err := middleware.GetClaimsFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid token"})
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid user ID in token"})
		return
	}

	var file models.File
	err = fc.FileService.GetFileByID(&file, fileID, uint(userID))
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	filePath := filepath.Join("uploads", file.FileName)

	fileToDownload, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer fileToDownload.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}
func (fc *FileController) FilesBuildSearchIndex(w http.ResponseWriter, r *http.Request) {
	fs := fc.FileService
	files, err := fs.AllFiles()
	if err != nil {
		log.Println("Error getting all files: ", err)
		http.Error(w, "Error getting files", http.StatusInternalServerError)
		return
	}

	if files == nil || len(files) == 0 {
		log.Println("No files found to index")
		http.Error(w, "No files found", http.StatusNotFound)
		return
	}

	log.Printf("Found %d files to index", len(files))

	for _, file := range files {

		log.Printf("Indexing file: %s", file.FileName)

		err := file.AddToIndex()
		if err != nil {
			log.Printf("Error adding file %s to index: %s", file.FileName, err)

		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Files indexed successfully"))
}

func (fc *FileController) SearchFiles(w http.ResponseWriter, r *http.Request) {
	var searchQuery struct {
		Query string `json:"query"`
	}

	var files []models.File
	var err error

	fs := fc.FileService

	if err := json.NewDecoder(r.Body).Decode(&searchQuery); err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid input"})
		return
	}

	if searchQuery.Query != "" {
		ids := models.FileSearch(searchQuery.Query)
		files, err = fs.FindFiles(ids)
		if err != nil {
			log.Println("Error getting files: ", err)
			http.Error(w, "Error getting files", http.StatusInternalServerError)
			return
		}
	} else {
		files, err = fs.AllFiles()
		if err != nil {
			log.Println("Error getting all files: ", err)
			http.Error(w, "Error getting files", http.StatusInternalServerError)
			return
		}
	}

	filesView := views.FilesView{
		Files: files,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filesView)
}
