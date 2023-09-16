package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mariogarzac/tecpool/pkg/config"
	"github.com/mariogarzac/tecpool/pkg/db"
	"github.com/mariogarzac/tecpool/pkg/models"
	"github.com/mariogarzac/tecpool/pkg/render"
)

// Creates a repo with the app configuration passed from main
var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Renders the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	var stringMap = map[string]string{}

	isLoggedIn := m.App.Session.GetBool(r.Context(), "isLoggedIn")
	// Check if user is logged in
	if isLoggedIn {

		stringMap["name"] = m.App.Session.GetString(r.Context(), "name")

		// Renders user dashboard
		render.RenderTemplate(w, r, "dashboard.page.html", &models.TemplateData{
			StringMap:  stringMap,
			IsLoggedIn: isLoggedIn,
		})

	} else {
		// Renders homepage with login or register buttons
		render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{
			IsLoggedIn: isLoggedIn,
		})
		log.Println("User not logged in")
	}
}

func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "register.page.html", &models.TemplateData{})
}

func (m *Repository) PostRegister(w http.ResponseWriter, r *http.Request) {

	// get values from the form
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	pass := r.FormValue("password")
	email := r.FormValue("email")
	phone := r.FormValue("phone_number")
	dob := r.FormValue("dob")

	// try to add the user to the database and return an error if it fails
	err := db.RegisterUser(fname, lname, pass, phone, email, dob)
	stringMap := map[string]string{}

	if err != nil {
		// render the register template with an error message
		stringMap["error_msg"] = "Error creating your account"
		render.RenderTemplate(w, r, "register.page.html", &models.TemplateData{
			StringMap: stringMap,
		})
	} else {
		// Render the login page on a success
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// func (m *Repository)Dashboard(w http.ResponseWriter, r *http.Request) {
//     render.RenderTemplate(w, r, "dashboard.page.html", &models.TemplateData{})
// }

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	// check if the user is already logged in if they are, it will display the Login
	// if they are not, the login form will be presented
	isLogged := m.App.Session.GetBool(r.Context(), "isLoggedIn")
	if isLogged {
		log.Println("Redirecting to dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

	} else {
		render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{})
	}
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {

	// get values from form
	email := r.FormValue("email")
	password := r.FormValue("password")

	// check if the user exits in the db
	err := db.ValidateUserInfo(email, password)

	stringMap := map[string]string{}

	if err != nil {
		// log the error
		log.Println("Wrong username or password", err)

		// set an error message for if the username or password are wrong
		stringMap["error_msg"] = "Wrong username or password"

		// render the template with the error message
		render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{
			StringMap: stringMap,
		})
	} else {
		// get the user's name by their email
		name, err := db.GetNameByEmail(email)
		userId, err := db.GetUserId(email)

		if err != nil {
			// log the error
			log.Println("Error getting the user's name ", err)

			// set an error message for if the username or password are wrong
			stringMap["error_msg"] = "An error occured please try again"

			// render the template with the error message
			render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{
				StringMap: stringMap,
			})
			return
		}

		// set cookie to logged in
		m.App.Session.Put(r.Context(), "isLoggedIn", true)
		m.App.Session.Put(r.Context(), "name", name)
		m.App.Session.Put(r.Context(), "userId", userId)

		// save cookie to db

		// if  err != nil {
		//     log.Println("Error saving session to db ", err)
		//     return
		// }

		// redirect to the dashboard
		log.Println("Log in success")
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (m *Repository) CreateTrip(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "dashboard.page.html", &models.TemplateData{})
}

func (m *Repository) PostCreateTrip(w http.ResponseWriter, r *http.Request) {
	carModel := r.FormValue("car_model")
	licensePlate := r.FormValue("plate")
	departureTime := r.FormValue("departure_time")
	userId := m.App.Session.GetInt(r.Context(), "userId")

	db.CreateTrip(carModel, licensePlate, departureTime, userId)

	// render.RenderTemplate(w, r, "create-trip.page.html", &models.TemplateData{})
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Fetch the 4 most recent trips from the database
	rows, err := db.GetRecentTrips()
	if err != nil {
		log.Println("Error fetching recent trips:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
			http.Error(w, "Internal server error", http.StatusInternalServerError)
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

	// Check for errors after the rows.Next() loop
	if err = rows.Err(); err != nil {
		log.Println("Error processing rows:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Render the dashboard template with the trips data
	render.RenderTemplate(w, r, "dashboard.page.html", &models.TemplateData{
		Trips: trips,
	})
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
