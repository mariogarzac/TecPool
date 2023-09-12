package config

import (
	"log"
	"html/template"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
    UseCache bool
    TemplateCache map[string]*template.Template
    InfoLog *log.Logger
    Session *scs.SessionManager
}
