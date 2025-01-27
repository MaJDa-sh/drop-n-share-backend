package views

import (
	"drop-n-share/internal/models"
	"log"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Conn   *websocket.Conn
	roomID string
}

func NewWebSocketHandler(conn *websocket.Conn, roomID string) *WebSocketHandler {
	return &WebSocketHandler{
		Conn:   conn,
		roomID: roomID,
	}
}

func (ws *WebSocketHandler) ReceiveMessage() (string, error) {
	_, msg, err := ws.Conn.ReadMessage()
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

func (ws *WebSocketHandler) GetRoomID() string {
	return ws.roomID
}

func (ws *WebSocketHandler) WriteMessage(messageType int, data []byte) error {
	return ws.Conn.WriteMessage(messageType, data)
}

func (ws *WebSocketHandler) Close() {
	err := ws.Conn.Close()
	if err != nil {
		log.Println("Error closing WebSocket connection:", err)
	}
}

func (ws *WebSocketHandler) Register() {
	models.RegisterClient(ws.Conn)
	models.RegisterRoomClient(ws.roomID, ws.Conn)
}

func (ws *WebSocketHandler) Unregister() {
	models.RemoveRoomClient(ws.roomID, ws.Conn)
	models.RemoveClient(ws.Conn)
}

func (ws *WebSocketHandler) BroadcastMessage(message []byte) {
	if clients, exists := models.RoomClients[ws.roomID]; exists {
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error sending message to client:", err)
			}
		}
	}
}
