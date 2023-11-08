package handlers

import (
	"log"
	"net/http"

	"github.com/mariogarzac/tecpool/pkg/config"
	"github.com/mariogarzac/tecpool/pkg/db"
	"github.com/mariogarzac/tecpool/pkg/models"
	"github.com/mariogarzac/tecpool/pkg/render"
)

// Creates a repo with the app configuration passed from main
var Repo *Repository

type Repository struct {
	App *config.AppConfig
    Server *Server
}

func NewRepo(a *config.AppConfig, s *Server) *Repository {
	return &Repository{
		App: a,
        Server: s,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Renders the home page
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


	isLoggedIn := m.App.Session.GetBool(r.Context(), "isLoggedIn")
	// Check if user is logged in
	if isLoggedIn {

		// Renders user dashboard with their name
        if err != nil {
            log.Println("Error getting user's name: ", err)
            return
        }

		render.RenderTemplate(w, r, "dashboard.page.html", &models.TemplateData{
            Trips: trips,
		})

	} else {
		// Renders login page
		render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{
	})
	}
}
