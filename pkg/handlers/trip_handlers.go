package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"strconv"

	"github.com/go-chi/chi"
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

    // Add the user who created the trip to the trip
    tripId, err := db.GetTripID()
    if err != nil {
        log.Println("Error getting tripId: ", err)
        return 
    }

    db.AddUserToTrip(uint64(userId), uint64(tripId))

    // create a chat room with that trip id
    m.CreateTripRoom(w, r)
    
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

func (m *Repository) JoinTrip(w http.ResponseWriter, r *http.Request){

    // Get id from the url and turn cast it into an int
    id := chi.URLParam(r, "tripId")
    tripId,err := strconv.Atoi(id)

    if err != nil {
        log.Println("Error getting tripId ", err)
    }

    // Get user info from cookies
    userId := m.App.Session.GetInt(r.Context(), "userId")

    if err != nil {
        log.Println("Error getting name: ", err)
        return
    }

    // Check if the trip exists
    if tripExists := db.DoesTripExist(tripId); !tripExists {
        log.Println("Room does not exist")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
        // show error message
    }

    // Check if user is already in the trip
    if inTrip := db.IsUserInTrip(tripId, userId); inTrip {
        log.Println("User is already in the room")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
        // show error message
    }
    
    // Add user to trip
    err = db.AddUserToTrip(uint64(userId), uint64(tripId));
    
    if err != nil {
        log.Println("Error adding user to trip: ", err)
        return
        // show error message
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) ActiveTrips(w http.ResponseWriter, r *http.Request){

    userId := m.App.Session.GetInt(r.Context(), "userId")

    userTrips := make(map[int]*models.Trip)

    // Get user trips from the db
    userTrips, err := db.GetUserTrips(userId)
    if err != nil {
        log.Println("Error getting user trips: ", err)
        return 
    }

    // Render trips in the active trips page
    render.RenderTemplate(w, r, "trips.page.html", &models.TemplateData{
        UserTrips: userTrips,
    })
}
