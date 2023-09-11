package render

import (
    "bytes"
    "log"
    "net/http"
    "path/filepath"
    "text/template"

    "github.com/mariogarzac/tecpool/pkg/config"
    "github.com/mariogarzac/tecpool/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig){
    app = a
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error{

    // New dictionary that will hold the template data mapping it to its name
    var templateCache map[string]*template.Template

    // maps the name of the template to its data
    templateCache,_ = CreateTemplateCache()

    // try to get the template data from the cache
    t, ok := templateCache[tmpl]
    if !ok {
        log.Fatal("Could not get template from template cache")
    }

    // executes the template data before running to check for any errors
    buf := new(bytes.Buffer)
    err := t.Execute(buf, td)
    if err != nil {
        log.Println("Error writing template to browser", err)
    }
    return err
}

func CreateTemplateCache() (map[string]*template.Template, error){
    cache := map[string]*template.Template{}

    pages, err := filepath.Glob("./templates/*.page.html")

    if err != nil {
        return cache , err
    }

    for _, page := range pages {
        filename := filepath.Base(page)
        ts, err := template.New(filename).ParseFiles(page)

        if err != nil {
            return cache , err
        }

        matches, err := filepath.Glob("./templates/*.layout.html")

        if len(matches) > 0 {
            ts, err = template.New(filename).ParseFiles(page)

            if err != nil {
                return cache ,err
            }

        }
        cache[filename] = ts
    }

    return cache, nil
}
