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
    var stringMap = map[string]string{}

    isLogged := m.App.Session.GetBool(r.Context(), "is_logged_in")
    if isLogged {

        stringMap["name"] = m.App.Session.GetString(r.Context(), "name")

        render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{
            StringMap: stringMap,
        })

    }else{
        stringMap["name"] = "Register or Login"
        render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{
            StringMap: stringMap,
        })
        log.Println("User not logged in")
    }
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
    // check if the user is already logged in
    isLogged := m.App.Session.GetBool(r.Context(), "is_logged_in")
    if isLogged {
        log.Println("Redirecting to dashboard")
        http.Redirect(w, r, "/", http.StatusSeeOther)

    }else{
        render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{})
    }
}

func (m *Repository)PostLogin(w http.ResponseWriter, r *http.Request) {

    // get values from form
    email := r.FormValue("email")
    password := r.FormValue("password")

    // check if the user exits in the db
    err := db.ValidateUserInfo(email, password)

    stringMap := map[string]string{}

    if err != nil{
        // log the error
        log.Println("Wrong username or password", err)

        // set an error message for if the username or password are wrong
        stringMap["error_msg"] = "Wrong username or password"

        // render the template with the error message
        render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{
            StringMap: stringMap,
        })
    }else{
        // get the user's name by their email
        name, err := db.GetNameByEmail(email)

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
        m.App.Session.Put(r.Context(), "is_logged_in", true)
        m.App.Session.Put(r.Context(), "name", name)

        // save cookie to db

        if  err != nil {
            log.Println("Error saving session to db ", err)
            return
        }

        // redirect to the dashboard
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func (m *Repository)CreateTrip(w http.ResponseWriter, r *http.Request){

}
