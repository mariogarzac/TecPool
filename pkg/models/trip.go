package models

type Trip struct {
	TripID        int    `json:"tripId"`
	CarModel      string `json:"carModel"`
	Date          string `json:"date"`
	Time          string `json:"time"`
	UserID        int    `json:"userId"`
	LicensePlate  string `json:"licensePlate"`
	ChatName      string `json:"chatName"`
	Image         string `json:"image"`
	StartLocation string `json:"startLocation"`
}
