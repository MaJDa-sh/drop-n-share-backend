package main

import (
	"drop-n-share/internal/controllers"
	"drop-n-share/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"

	_ "github.com/lib/pq"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/users", controllers.WebSocketUserRoute)
	router.HandleFunc("/room/{id}", controllers.WebSocketRoomRoute)

	db := sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))

	routes.SetupRoutes(db, router)

	var version string
	err := db.QueryRow("select version()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB: ", version)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Println("Server started on :8000")
	if err := http.ListenAndServe(":8000", handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
