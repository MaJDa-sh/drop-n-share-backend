package routes

import (
	"drop-n-share/internal/controllers"
	"drop-n-share/internal/middleware"
	"drop-n-share/internal/services"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func SetupRoutes(db *sqlx.DB, router *mux.Router) {

	router.HandleFunc("/", controllers.HandleConnections)

	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	router.HandleFunc("/sign_in", userController.HandleSignIn).Methods("POST")

	router.HandleFunc("/sign_up", userController.HandleSignUp).Methods("POST")

	protected := router.PathPrefix("/protected").Subrouter()
	protected.Use(middleware.JWTMiddleware)

	protected.HandleFunc("/profile", userController.HandleGetUserByJWT).Methods("GET")

	fileService := services.NewFileService(db)
	fileController := controllers.NewFileController(fileService)

	protected.HandleFunc("/upload", fileController.HandleFileUpload).Methods("POST")
	protected.HandleFunc("/download", fileController.HandleFileDownload).Methods("GET")
	router.HandleFunc("/files/index", fileController.FilesBuildSearchIndex).Methods("GET")
	router.HandleFunc("/files/search", fileController.SearchFiles).Methods("POST")

	// protected.HandleFunc("/search", controllers.SearchFilesHandler)

}
