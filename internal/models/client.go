package models

import (
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

var Client = make(map[*websocket.Conn]bool)

var RoomClients = make(map[string]map[*websocket.Conn]bool)

func RegisterClient(conn *websocket.Conn) {
	Client[conn] = true
}

func RemoveClient(conn *websocket.Conn) {
	delete(Client, conn)
}

func GetClientList() []string {
	clientList := []string{}
	for client := range Client {
		clientList = append(clientList, client.RemoteAddr().String())
	}
	return clientList
}

func RegisterRoomClient(roomID string, conn *websocket.Conn) {

	if RoomClients[roomID] == nil {
		RoomClients[roomID] = make(map[*websocket.Conn]bool)
	}
	RoomClients[roomID][conn] = true
}

func RemoveRoomClient(roomID string, conn *websocket.Conn) {
	delete(RoomClients[roomID], conn)

	if len(RoomClients[roomID]) == 0 {
		delete(RoomClients, roomID)
	}
}

func GetRoomClientList(roomID string) []string {
	clientList := []string{}

	if clients, exists := RoomClients[roomID]; exists {
		for client := range clients {
			clientList = append(clientList, client.RemoteAddr().String())
		}
	}
	log.Println(len(strings.Join(clientList, "\n")))
	return clientList
}

func GetRoomClientCount(roomID string) int {
	if clients, exists := RoomClients[roomID]; exists {
		return len(clients)
	}
	return 0
}
