package main

import (
	"log"

	"github.com/mariogarzac/tecpool/pkg/config"
	"github.com/mariogarzac/tecpool/pkg/handlers"
	"github.com/mariogarzac/tecpool/pkg/render"
)


var portNumber = ":8080"
var app config.AppConfig

func main(){

    tc, err := render.CreateTemplateCache()
    if err != nil {
        log.Println("Cannot create template cache")
    }

    // used to test if templates are being rendered 
    app.TemplateCache = tc
    app.UseCache = false

    repo := handlers.NewRepo(&app)
    handlers.NewHandlers(repo)

    render.NewTemplates(&app)

    serveAndRoute()


}
