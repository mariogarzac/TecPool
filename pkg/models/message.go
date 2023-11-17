package models

type Message struct {
    Message string `json:"message"`
    UserID int `json:"userId"`
    ChatID int `json:"chatId"`
    ChatName string `json:"chatName"`
    Time []uint8 `json:"time"`
    Self bool `json:"self"`
}
