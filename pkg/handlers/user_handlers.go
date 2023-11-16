package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mariogarzac/tecpool/pkg/db"
	"github.com/mariogarzac/tecpool/pkg/models"
	"github.com/mariogarzac/tecpool/pkg/render"
)

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

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
    // check if the user is already logged in if they are, it will display the Login
    // if they are not, the login form will be presented
    isLogged := m.App.Session.GetBool(r.Context(), "isLoggedIn")
    if isLogged {
        log.Println("Redirecting to dashboard")
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    } else {
        render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{})
    }
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {

    // get values from form
    email := r.FormValue("email")
    password := r.FormValue("password")


    // check if the user exits in the db
    stringMap := map[string]string{}
    err := db.ValidateUserInfo(email, password)

    if err != nil {
        // log the error
        log.Println("Wrong username or password", err)

        // set an error message for if the username or password are wrong
        stringMap["error_msg"] = "Wrong username or password"

        // render the template with the error message
        render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{
            StringMap: stringMap,
        })
    } 

    // get the user's name by their email
    userId, err := db.GetUserIDByEmail(email)

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


    name,err := db.GetNameByID(userId)
    if err != nil {
        log.Println("Error getting name by id: ", err)
        return 
    }
    // set cookie to logged in
    m.App.Session.Put(r.Context(), "isLoggedIn", true)
    m.App.Session.Put(r.Context(), "userId", userId)
    m.App.Session.Put(r.Context(), "name", name)

    // redirect to the dashboard
    log.Println("Log in success")
    http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (m *Repository)ShowSettings(w http.ResponseWriter, r *http.Request){

    // Get User history
    userID := m.App.Session.GetInt(r.Context(), "userId")
    userTripHistory, err := db.GetUserTrips(userID)
    if err != nil {
        log.Println(err)
        return 
    }

    user := db.GetUserInfo(userID)

    render.RenderTemplate(w, r, "settings.page.html", &models.TemplateData{
        UserTrips: userTripHistory,
        Users: user,
    })
}

func (m *Repository)Logout(w http.ResponseWriter, r *http.Request){
    m.App.Session.Destroy(r.Context())

    render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{})
}

type PasswordUpdate struct {
    OldPassword string `json:"old_password"`
    NewPassword string `json:"new_password"`
}

func (m *Repository)ChangePassword(w http.ResponseWriter, r *http.Request){

    var update PasswordUpdate

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&update); err != nil{ 
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    userID := m.App.Session.GetInt(r.Context(), "userId")
    oldPassword := update.OldPassword
    newPassword := update.NewPassword

    log.Println(oldPassword, newPassword)

    err := db.ChangePassword(userID, oldPassword, newPassword)

    if err != nil {
        return 
    }

}
