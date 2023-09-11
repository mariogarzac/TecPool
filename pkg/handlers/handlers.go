package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mariogarzac/tecpool/pkg/config"
	// "github.com/mariogarzac/tecpool/pkg/models"
	// "github.com/mariogarzac/tecpool/pkg/render"
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

// func (m *Repository)Login(w http.ResponseWriter, r *http.Request) {
//     render.RenderTemplate(w,r, "login.page.html", &models.TemplateData{})
// }

func (m *Repository)Login(c echo.Context) error {
    return c.String(http.StatusOK, "hello from handlers")
}
