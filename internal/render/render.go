package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/dorasaicu12/booking/internal/config"
	"github.com/dorasaicu12/booking/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate" : HumanDate,
}

var app *config.AppConfig
var pathToTemplate ="./templates"

// NewTemplates sets the config for the template package
func NewRendered(a *config.AppConfig) {
	app = a
}
//return time in yyy mmmm dddd
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func AddDefaultData(td *models.TemplateData,r *http.Request) *models.TemplateData {
	td.Flash =app.Session.PopString(r.Context(),"flash")
	td.Error =app.Session.PopString(r.Context(),"error")
	td.Warning =app.Session.PopString(r.Context(),"warning")
    td.CSRFToken = nosurf.Token(r)
	td.UrlGlobal="http://localhost:8080"
	if app.Session.Exists(r.Context(),"user_id"){
		td.IsAuthenticated=1
	}
	return td
}

// RenderTemplate renders a template
func Template(w http.ResponseWriter,r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
		return errors.New("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td,r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}
	return nil

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl",pathToTemplate))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl",pathToTemplate))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl",pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
