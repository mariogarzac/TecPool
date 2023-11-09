package handlers

import (
    "io"
    "log"
    "net/http"
    "strconv"

    "github.com/go-chi/chi"
    "github.com/mariogarzac/tecpool/pkg/models"
    "github.com/mariogarzac/tecpool/pkg/render"
    "golang.org/x/net/websocket"
)

type Hub struct {
    conns map[*websocket.Conn]bool
}

func NewHub() *Hub {
    return &Hub{
        conns: make(map[*websocket.Conn]bool),
    }
}

type Message struct {
    Message string `json:"message"`
    UserID string `json:"userId"`
    Self bool `json:"self"`
}

func (m *Repository)RenderChat(w http.ResponseWriter, r *http.Request) {

    stintMap := make(map[string]int)
    userId := chi.URLParam(r, "userId")

    uid, _ := strconv.Atoi(userId)

    tripId := chi.URLParam(r, "tripId")
    tid, _ := strconv.Atoi(tripId)

    stintMap["userId"] = uid
    stintMap["tripId"] = tid

    render.RenderTemplate(w, r, "chat.page.html", &models.TemplateData{
        StringIntMap: stintMap,
    })
}

func (s *Hub) HandleWs(ws *websocket.Conn) {
    log.Println("New connection from ", ws.RemoteAddr())

    s.conns[ws] = true
    s.readLoop(ws)
}

func (s *Hub)readLoop(ws *websocket.Conn) {

    var messageData map[string]string
    for {
        err := websocket.JSON.Receive(ws, &messageData)

        if err != nil {
            if err == io.EOF {
                // The client closed the connection
                log.Println("Client closed the connection")
                return
            } else {
                return
            }
        }

        msg := messageData["message"]
        userId := messageData["userId"]

        s.broadcast(msg, ws, userId)
    }
}

func (s *Hub)broadcast(b string, sender *websocket.Conn, userId string) {
    log.Println(b)
    for ws := range s.conns {
        jsonMsg := &Message{
            Message: b,
            UserID: userId,
            Self: ws == sender,
        }

        encodedMsg, _, err := websocket.JSON.Marshal(jsonMsg)
        if err != nil {
            log.Printf("Error encoding JSON message: %v", err)
            return
        }

        go func(ws *websocket.Conn) {
            // send the json message to the front end here
            if _, err := ws.Write(encodedMsg); err != nil {
                return 
            }
        }(ws)
    }
}
