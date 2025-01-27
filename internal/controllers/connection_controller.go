package controllers

import (
	"drop-n-share/internal/models"
	"drop-n-share/internal/views"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	models.RegisterClient(conn)
	log.Println("New WebSocket connection")

	views.SendClientListToClients()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			models.RemoveClient(conn)
			views.SendClientListToClients() 
			break
		}
		log.Printf("Received message: %s", message)

		views.BroadcastMessage(message)
	}
}
