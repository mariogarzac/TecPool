package models

// Makes it easier to pass maps around
type TemplateData struct {
    StringMap map[string]string
    IntMap map[int]int
    FloatMap map[float32]float32
    Data map[string]interface{}
}
