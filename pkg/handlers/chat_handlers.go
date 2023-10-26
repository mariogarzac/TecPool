package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mariogarzac/tecpool/pkg/db"
	"github.com/mariogarzac/tecpool/pkg/models"
	"github.com/mariogarzac/tecpool/pkg/render"
	"golang.org/x/net/websocket"
)

type User struct{
    UserID int `json:"userID"`
    Name string `json:"name"`
    Conn *websocket.Conn
}

type Rooms struct {
    RoomID int `json:"roomID"`
    Users map[int]*User
    Broadcast chan []byte
}

func (m *Repository) CreateRoom(w http.ResponseWriter, r *http.Request) {

    // Get name and email from cookies
    name := m.App.Session.GetString(r.Context(), "name")
    email := m.App.Session.GetString(r.Context(), "email")

    // Get tripId from DB
    tripId, err := db.GetTripID()
    if err != nil {
        log.Fatal("Error getting tripId: ", err)
    }

    // Get userId from DB
    userId, err := db.GetUserId(email)

    if err != nil {
        log.Fatal("Error getting userId: ", err)
    }

    // Create User and Room to add to the TripMap
    user := &User{
        Name: name,
        UserID: userId,
    }

    room := &Rooms{
        RoomID: tripId,
        Users: map[int]*User{
            userId: user,
        },
    }

    m.TripMap[tripId] = room
}

func (m *Repository) JoinTrip(w http.ResponseWriter, r *http.Request){

    // Get id from the url and turn cast it into an int
    id := chi.URLParam(r, "tripId")
    tripId,err := strconv.Atoi(id)

    if err != nil {
        log.Println("Error getting tripId ", err)
    }

    // Get user info from cookies
    // TODO: Should change to userId
    name := m.App.Session.GetString(r.Context(), "name")
    email := m.App.Session.GetString(r.Context(), "email")

    if err != nil {
        log.Fatal("Error getting tripId: ", err)
    }

    // Get userId from email
    userId, err := db.GetUserId(email)
    if err != nil {
        log.Println("Error getting userId: ", err)
        return
    }

    // Check if the trip exists
    if _, exists := m.TripMap[tripId]; !exists {
        log.Println("Room does not exist")
        return
    }

    // Check if user is already in the trip
    if _, inRoom := m.TripMap[tripId].Users[userId]; inRoom {
        log.Println("User is already in the room")
        return
    }
    
    // Add user to trip and redirect to dashboard
    // TODO: Add to database as well
    m.TripMap[tripId].Users[userId] = &User{
        Name:   name,
        UserID: userId,
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) ActiveTrips(w http.ResponseWriter, r *http.Request){

    render.RenderTemplate(w, r, "trips.page.html", &models.TemplateData{})
}

func (m *Repository) printTrips(){
    for tripID, room := range m.TripMap {
        fmt.Printf("Trip ID: %d\n", tripID)

        for userID, user := range room.Users {
            fmt.Printf("User ID: %d\n", userID)
            fmt.Printf("User Name: %s\n", user.Name)
        }
        fmt.Println() 
    }
}

// func (m *Repository) JoinTrip(ws *websocket.Conn){
//     defer ws.Close()
//
//     log.Println("Comming from ", ws.RemoteAddr())
//     tripId, err := strconv.Atoi(ws.Request().URL.Query().Get("tripId"))
//    
//     if err != nil {
//         log.Println("Error getting trip id ", err)
//         return
//     }
//
//     log.Println("Trip id is ", tripId)
//     //
//     // if err := websocket.Message.Send(ws, "You have successfully joined the chat."); err != nil {
//     //     log.Println("Error sending success message:", err)
//     // }
//
// }

// func (m *Repository) readLoop(ws *websocket.Conn){
//     buf := make([]byte, 1024)
//
//     for {
//         n, err := ws.Read(buf)
//
//         if err != nil {
//             if err == io.EOF{
//                 fmt.Println("Connection closed...")
//                 break
//             }
//             fmt.Println("read error: ", err)
//             continue
//         }
//
//         msg := buf[:n]
//         s.broadcast(msg)
//     }
// }

// func (m *Repository) broadcast(b []byte){
// need the room id here
// tripID := 
// for ws := range m.TripMap[tripId] {
//     go func(ws *websocket.Conn){
//         if _, err := ws.Write(b); err != nil {
//             log.Println("write error: ", err)
//         }
//     }(ws)
// }
// }
