package handlers

import (
    "encoding/json"
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
    conns map[*websocket.Conn]int
}

func NewHub() *Hub {
    return &Hub{
        conns: make(map[*websocket.Conn]int),
    }
}

type NewChatNameRequest struct{
    ChatName string `json:"chatName"`
    ChatID string `json:"chatId"`
}

func (m *Repository)UpdateChatName(w http.ResponseWriter, r *http.Request){
    var request NewChatNameRequest

    // Decode the JSON request body into the struct
    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        log.Println(err)
        http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
        return
    }

    // Access the new group name using request.NewGroupName
    // Perform the update in the database or handle as needed
    newChatName := request.ChatName
    chatId := request.ChatID
    log.Println("Received ", newChatName, chatId)

    cid, err := strconv.Atoi(chatId)
    if err != nil {
        log.Println("Error getting chat id", err)
        return 
    }

    if err := db.UpdateGroupName(newChatName, cid); err != nil{
        log.Println("Could not update group name")
        return
    }

    // Respond with a success message (optional)
    response := map[string]string{"newChatName": newChatName, "message": "Group name updated successfully"}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)

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

    // Fetch the user's active trips
    userTrips, err := db.GetUserTrips(uid)
    if err != nil {
        log.Println("Error getting user trips: ", err)
        return
    }

    chatName, err := db.GetTripTitle(tid)
    if err != nil {
        log.Println("Error getting chat name", err)
        return 
    }

    tripParticipants,err := db.GetParticipants(tid)

    if err != nil {
        log.Println("Error getting trip participants", err)
        return 
    }

    render.RenderTemplate(w, r, "chat.page.html", &models.TemplateData{
        StringIntMap: stintMap,
        Messages:     messages,
        UserTrips:    userTrips, 
        ChatName:     chatName,
        Users:        tripParticipants,
    })
}

func (s *Hub) HandleWs(ws *websocket.Conn) {
    log.Println("New connection from ", ws.RemoteAddr())

    var messageData map[string]string
    err := websocket.JSON.Receive(ws, &messageData)

    if err != nil {
        if err == io.EOF {
            return
        } 
    }

    tid,err := strconv.Atoi(messageData["tripId"])

    if err != nil {
        return 
    }

    s.conns[ws] = tid
    s.readLoop(ws)
}

func (s *Hub)readLoop(ws *websocket.Conn) {

    var messageData map[string]string
    for {
        err := websocket.JSON.Receive(ws, &messageData)

        if err != nil {
            if err == io.EOF {
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
    for ws,tid := range s.conns {
        if tid == tripId{
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

}
