package models

import "github.com/mariogarzac/tecpool/pkg/forms"

// Makes it easier to pass maps around
type TemplateData struct {
    StringMap map[string]string
    IntMap map[int]int
    FloatMap map[float32]float32
    Data map[string]interface{}
    Forms *forms.Form
    IsLoggedIn bool
    Trips []map[string]interface{}
    UserTrips map[int]*Trip
}
