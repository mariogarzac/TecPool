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

type Server struct {
   conns map[*websocket.Conn]bool
}

func NewServer() *Server {
    return &Server{
        conns: make(map[*websocket.Conn]bool),
    }
}

func (m *Repository)RenderChat(w http.ResponseWriter, r *http.Request) {

    stintMap := make(map[string]int)
    userId := chi.URLParam(r, "userId")

    uid, _ := strconv.Atoi(userId)
    
    tripId := chi.URLParam(r, "tripId")
    tid, _ := strconv.Atoi(tripId)

    stintMap["userId"] = uid
    stintMap["tripId"] = tid

    log.Println(userId, tripId)

    render.RenderTemplate(w, r, "chat.page.html", &models.TemplateData{
         StringIntMap: stintMap,
    })
}

func (s *Server) HandleWs(ws *websocket.Conn) {
    log.Println("New connection from ", ws.RemoteAddr())

    s.conns[ws] = true

    s.readLoop(ws)
}

func (s *Server)readLoop(ws *websocket.Conn) {
    buf := make([]byte, 1024)
    for {
        n, err := ws.Read(buf)

        if err != nil {
            if err == io.ErrUnexpectedEOF {
                break
            }
            log.Println(err)
            return 
        }

        msg := buf[:n]

        s.broadcast(msg)
    }
}

func (s *Server)broadcast(b []byte) {
    for ws := range s.conns {
        go func(ws *websocket.Conn) {
            if _, err := ws.Write(b); err != nil {
                return 
            }
        }(ws)
    }
}
