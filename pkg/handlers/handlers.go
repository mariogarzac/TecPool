package handlers

import (
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
}

func NewRepo(a *config.AppConfig) *Repository{
    return &Repository{
        App: a,
    }
}

func NewHandlers(r *Repository){
    Repo = r
}

func (m *Repository)Home(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

func (m *Repository)Register(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, r, "register.page.html", &models.TemplateData{})
}

func (m *Repository)PostRegister(w http.ResponseWriter, r *http.Request) {

    // get values from the form
    fname := r.FormValue("fname")
    lname := r.FormValue("lname")
    pass := r.FormValue("password")
    email := r.FormValue("email")
    phone := r.FormValue("phone_number")
    dob := r.FormValue("dob")

    // try to add the user to the database and return an error if it fails
    err := db.RegisterUser(fname, lname, pass, phone, email, dob)

    if err != nil {
        // render the register template with an error message
        stringMap := map[string]string{}
        stringMap["error_msg"] = "Error creating your account"
        render.RenderTemplate(w, r, "register.page.html", &models.TemplateData{
            StringMap: stringMap,
        })
    }else{
        // Render the login page on a success
        http.Redirect(w, r, "/login", http.StatusSeeOther)
    }

}

func (m *Repository)Login(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{})
}

func (m *Repository)PostLogin(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    password := r.FormValue("password")
    err := db.ValidateUserInfo(email, password)

    if err != nil{
        stringMap := map[string]string{}
        stringMap["error_msg"] = "Wrong username or password"
        render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{
            StringMap: stringMap,
        })
    }else{
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}
