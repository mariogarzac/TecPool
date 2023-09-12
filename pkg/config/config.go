package config

import (
	"database/sql"
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
    UseCache bool
    TemplateCache map[string]*template.Template
    InfoLog *log.Logger
    Session *scs.SessionManager
    DataBase *sql.DB
}
