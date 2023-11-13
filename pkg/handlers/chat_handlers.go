package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mariogarzac/tecpool/pkg/db"
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


func (m *Repository)RenderChat(w http.ResponseWriter, r *http.Request) {

    stintMap := make(map[string]int)

    userId := chi.URLParam(r, "userId")
    uid, _ := strconv.Atoi(userId)

    tripId := chi.URLParam(r, "tripId")
    tid, _ := strconv.Atoi(tripId)

    stintMap["tripId"] = tid
    stintMap["userId"] = uid

    tid,_ = strconv.Atoi(tripId)
    messages,_ := db.LoadMessages(tid, uid)

    for _, msg := range messages {
        fmt.Println(msg.Self)
    }


    render.RenderTemplate(w, r, "chat.page.html", &models.TemplateData{
        StringIntMap: stintMap,
        Messages: messages,
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
        tripId := messageData["tripId"]

        uid, err := strconv.Atoi(userId)
        if err != nil {
            log.Println(err)
            return 
        }

        tid, err := strconv.Atoi(tripId)
        if err != nil {
            log.Println(err)
            return 
        }

        s.broadcast(msg, ws, tid, uid)
    }
}

func (s *Hub)broadcast(msg string, sender *websocket.Conn, tripId, userId int) {
    for ws := range s.conns {
        jsonMsg := &models.Message{
            Message: msg,
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

    db.SaveMessage(tripId, userId, msg)

}
