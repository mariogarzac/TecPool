package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mariogarzac/tecpool/pkg/config"
	"github.com/mariogarzac/tecpool/pkg/handlers"
	"github.com/mariogarzac/tecpool/pkg/render"
)


const portNumber = ":8080"
var app config.AppConfig

func main(){

    app.UseCache = true
    app.Session = scs.New()
    app.Session.Lifetime = 24 * time.Hour
    app.Session.Cookie.Persist = true
    app.Session.Cookie.SameSite = http.SameSiteLaxMode


    tc, err := render.CreateTemplateCache()

    if err != nil {
        log.Fatal("Cannot create template cache: ", err)
    }

    // used to test if templates are being rendered 
    app.TemplateCache = tc

    repo := handlers.NewRepo(&app)
    handlers.NewHandlers(repo)

    render.NewTemplates(&app)

    srv := &http.Server {
        Addr: portNumber,
        Handler: routes(&app),
    }

    err = srv.ListenAndServe()

    if err != nil {
        log.Fatal(err)
    }

}
