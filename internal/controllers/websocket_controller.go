package controllers

import (
	"drop-n-share/internal/models"
	"drop-n-share/internal/views"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var hub = models.NewHub()

type WebSocketController struct {
	WebSocketView *views.WebSocketHandler
	Hub           *models.Hub
}

func NewWebSocketController(wsView *views.WebSocketHandler, hub *models.Hub) *WebSocketController {
	return &WebSocketController{
		WebSocketView: wsView,
		Hub:           hub,
	}
}

func (controller *WebSocketController) HandleClientMessages(listUsers bool) {
	room := controller.Hub.GetRoom(controller.WebSocketView.GetRoomID())
	room.RegisterClient(controller.WebSocketView.Conn)

	if listUsers {
		controller.notifyUserList()
	}

	for {
		msg, err := controller.WebSocketView.ReceiveMessage()
		if err != nil {
			log.Println("Error receiving message:", err)
			break
		}

		room.BroadcastMessage([]byte(msg))
	}

	room.RemoveClient(controller.WebSocketView.Conn)

	if listUsers {
		controller.notifyUserList()
	}
}

func (controller *WebSocketController) notifyUserList() {
	roomID := controller.WebSocketView.GetRoomID()
	log.Println("Notifying user list for room:", roomID)

	room := controller.Hub.GetRoom(roomID)
	userList := room.GetClientList()

	currentUserID := controller.WebSocketView.Conn.RemoteAddr().String()

	response := map[string]interface{}{
		"current_user": currentUserID,
		"users":        userList,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling response to JSON:", err)
		return
	}

	log.Println("User list message:", string(responseJSON))

	room.BroadcastMessage(responseJSON)
}

func WebSocketUserRoute(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	wsHandler := views.NewWebSocketHandler(conn, "users")

	controller := NewWebSocketController(wsHandler, hub)
	controller.HandleClientMessages(true)

	defer hub.GetRoom("users").RemoveClient(conn)
}

func WebSocketRoomRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["id"]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	wsHandler := views.NewWebSocketHandler(conn, roomID)

	room := hub.GetRoom(roomID)

	room.RegisterClient(conn)

	controller := NewWebSocketController(wsHandler, hub)
	controller.HandleClientMessages(false)

	defer room.RemoveClient(conn)
}

func GetUserListHandler(w http.ResponseWriter, r *http.Request) {
	userList := models.GetClientList()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(userList); err != nil {
		http.Error(w, "Failed to encode user list", http.StatusInternalServerError)
		return
	}
}
