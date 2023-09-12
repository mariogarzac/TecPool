package handlers

import (
	"fmt"
	"net/http"

	"github.com/mariogarzac/tecpool/pkg/config"
	"github.com/mariogarzac/tecpool/pkg/models"
	"github.com/mariogarzac/tecpool/pkg/render"
)

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

func (m *Repository)Login(w http.ResponseWriter, r *http.Request) {
    render.RenderTemplate(w, r, "login.page.html", &models.TemplateData{})
}

func (m *Repository)PostLogin(w http.ResponseWriter, r *http.Request) {
    user := r.FormValue("username")
    pass := r.FormValue("password")
    w.Write([]byte(fmt.Sprintf("posted to login with user %s and password %s", user, pass)))
}
