package handlers

import (
	"fmt"
	"log"
	"net/http"

    "github.com/mariogarzac/tecpool/pkg/db"
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

// Creates the room for the trip. This adds the trip to memory
func (m *Repository) CreateTripRoom(w http.ResponseWriter, r *http.Request) {

    // Get name and email from cookies
    userId := m.App.Session.GetInt(r.Context(), "userId")

    name, err := db.GetNameByID(userId)
    if err != nil {
        log.Println("Error getting name ", err)
        return
    }

    // Get tripId from DB
    tripId, err := db.GetTripID()
    if err != nil {
        log.Fatal("Error getting tripId: ", err)
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
