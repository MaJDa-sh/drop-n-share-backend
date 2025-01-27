package views

import (
	"drop-n-share/internal/models"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

func BroadcastMessage(message []byte) {
	for client := range models.Client {
		go func(client *websocket.Conn) {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Error sending message:", err)
				client.Close()
				models.RemoveClient(client)
				SendClientListToClients()
			}
		}(client)
	}
}

func SendClientListToClients() {
	clientList := models.GetClientList()
	clientListMessage := joinClientList(clientList)

	for client := range models.Client {
		err := client.WriteMessage(websocket.TextMessage, []byte(clientListMessage))
		if err != nil {
			log.Println("Error sending client list:", err)
		}
	}
}

func joinClientList(clientList []string) string {
	return strings.Join(clientList, "\n")
}
