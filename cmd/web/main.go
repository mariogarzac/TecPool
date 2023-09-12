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

    // set app configuration for cookies, cache and sessions
    app.UseCache = true
    app.Session = scs.New()
    app.Session.Lifetime = 24 * time.Hour
    app.Session.Cookie.Persist = true
    app.Session.Cookie.SameSite = http.SameSiteLaxMode


    // render the templates
    tc, err := render.CreateTemplateCache()

    // check if templates were rendered correctly
    if err != nil {
        log.Fatal("Cannot create template cache: ", err)
    }

    // used to test if templates are being rendered 
    app.TemplateCache = tc

    // passes app configuration and templates to handlers
    repo := handlers.NewRepo(&app)
    handlers.NewHandlers(repo)

    render.NewTemplates(&app)

    // start the server with the port number and indicate where routes are
    srv := &http.Server {
        Addr: portNumber,
        Handler: routes(&app),
    }

    err = srv.ListenAndServe()

    if err != nil {
        log.Fatal(err)
    }

}
