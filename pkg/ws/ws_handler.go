package ws

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error upgrading", err)
        return 
    }

	roomID := chi.URLParam(r, "roomId")
	clientID := chi.URLParam(r, "userId")
	username := chi.URLParam(r, "username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

// type RoomRes struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// }


// type ClientRes struct {
// 	ID       string `json:"id"`
// 	Username string `json:"username"`
// }
//
// func (h *Handler) GetClients(c *gin.Context) {
// 	var clients []ClientRes
// 	roomId := c.Param("roomId")
//
// 	if _, ok := h.hub.Rooms[roomId]; !ok {
// 		clients = make([]ClientRes, 0)
// 		c.JSON(http.StatusOK, clients)
// 	}
//
// 	for _, c := range h.hub.Rooms[roomId].Clients {
// 		clients = append(clients, ClientRes{
// 			ID:       c.ID,
// 			Username: c.Username,
// 		})
// 	}
//
// 	c.JSON(http.StatusOK, clients)
// }
//
