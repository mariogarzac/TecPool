package render

import (
    "bytes"
    "log"
    "net/http"
    "path/filepath"
    "html/template"

    "github.com/mariogarzac/tecpool/pkg/config"
    "github.com/mariogarzac/tecpool/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig){
    app = a
}

func AddDefautlData(td *models.TemplateData) *models.TemplateData {
    return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {

    // New dictionary that will hold the template data mapping it to its name
    var templateCache map[string]*template.Template

    if app.UseCache {
        templateCache = app.TemplateCache
    }else{
        templateCache, _ = CreateTemplateCache()
    }

    // try to get the template data from the cache
    t, ok := templateCache[tmpl]
    if !ok {
        log.Fatal("Could not get template from template cache")
    }

    // executes the template data before running to check for any errors
    buf := new(bytes.Buffer)

    td = AddDefautlData(td)

    err := t.Execute(buf, td)
    if err != nil {
        log.Println("Error executing template ", err)
    }

    // render the template
    _, err = buf.WriteTo(w)
    if err != nil {
        log.Println("Error writing template to browser", err)
       
    }
}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}

