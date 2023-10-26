package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mariogarzac/tecpool/pkg/db"
	"github.com/mariogarzac/tecpool/pkg/models"
	"github.com/mariogarzac/tecpool/pkg/render"
)

func (m *Repository) CreateTrip(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "dashboard.page.html", &models.TemplateData{})
}

func (m *Repository) PostCreateTrip(w http.ResponseWriter, r *http.Request) {
	carModel := r.FormValue("car_model")
	licensePlate := r.FormValue("plate")
	departureTime := r.FormValue("departure_time")
	userId := m.App.Session.GetInt(r.Context(), "userId")

    err := db.CreateTrip(carModel, licensePlate, departureTime, userId)
    if err != nil {
        log.Fatal("Error creating trip", err)
    }

    // create a chat room with that trip id
    m.CreateRoom(w, r)
    
	// render.RenderTemplate(w, r, "create-trip.page.html", &models.TemplateData{})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) SearchTripsHandler(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, `{"error": "Invalid request format"}`, http.StatusBadRequest)
		return
	}

	departureTimeString, ok := requestData["departureTime"]
	if !ok || departureTimeString == "" {
		log.Println("Error: Departure time is empty.")
		http.Error(w, `{"error": "Please provide a valid departure time."}`, http.StatusBadRequest)
		return
	}

	departureTime, err := time.Parse("2006-01-02T15:04", departureTimeString)
	if err != nil {
		log.Println("Error parsing departure time:", err)
		http.Error(w, `{"error": "Invalid departure time format"}`, http.StatusBadRequest)
		return
	}

	log.Println("Departure time", departureTime)

	// Fetch trips from the database that match the given departure time
	rows, err := db.SearchTripsByDepartureTime(departureTime)
	if err != nil {
		log.Println("Error fetching trips by departure time:", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var trips []map[string]interface{}
	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Println("Error scanning row:", err)
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
		trip := make(map[string]interface{})
		for i, col := range columns {
			// Check if the value is a byte slice ([]byte) and convert it to a string
			if v, ok := values[i].([]byte); ok {
				trip[col] = string(v)
			} else {
				trip[col] = values[i]
			}
		}
		trips = append(trips, trip)
	}
	log.Println("Fetched trips:", trips)

	// Check for errors after the rows.Next() loop
	if err = rows.Err(); err != nil {
		log.Println("Error processing rows:", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(trips); err != nil {
		log.Println("Error encoding trips to JSON:", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
	}

}
