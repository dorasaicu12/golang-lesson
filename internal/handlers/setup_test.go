package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dorasaicu12/booking/internal/config"
	"github.com/dorasaicu12/booking/internal/driver"
	"github.com/dorasaicu12/booking/internal/models"
	"github.com/dorasaicu12/booking/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)
 var app config.AppConfig
 var session *scs.SessionManager
func getRoutes() http.Handler{
			// what i am going to put the session
			gob.Register(models.Reservation{})
			// change this to true when in production
			app.InProduction = false
		
			// set up the session
			session = scs.New()
			session.Lifetime = 24 * time.Hour
			session.Cookie.Persist = true
			session.Cookie.SameSite = http.SameSiteLaxMode
			session.Cookie.Secure = app.InProduction
		
			app.Session = session
			db,err :=driver.ConnectSQL("root:@(127.0.0.1:3306)/booking")

			infoLog :=log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime)
			app.InfoLog=infoLog
			errorLog :=log.New(os.Stdout,"ERROR\t",log.Ldate|log.Ltime|log.Lshortfile)
			app.ErrorLog=errorLog
		
			tc, err := CreateTestTemplateCache()
			if err != nil {
				log.Fatal("cannot create template cache")
				
			}
		
			app.TemplateCache = tc
			app.UseCache = true

			repo := NewRepo(&app,db)
			NewHandlers(repo)
		
			render.NewRendered(&app)

			mux := chi.NewRouter()

mux.Use(middleware.Recoverer)
// mux.Use(NoSurf)
mux.Use(SessionLoad)

mux.Get("/", Repo.Home)
mux.Get("/about", Repo.About)
mux.Get("/general-quate", Repo.General)
mux.Get("/major-suite", Repo.Major)

mux.Get("/search-avai", Repo.Search_Avai)
mux.Post("/search-avai", Repo.Handle_Search_Avai)
mux.Post("/search-avai-json", Repo.AvaibilityJSON)

mux.Get("/contact", Repo.Contact)
mux.Get("/make-reservation", Repo.Reservation)
mux.Post("/make-reservation", Repo.PostReservation)
mux.Get("/reservation-summery",Repo.ReservationSummery)

fileServer := http.FileServer(http.Dir("./static/"))
mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
return mux
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves session data for current request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

var pathToTemplate ="./templates"
var functions = template.FuncMap{}
func CreateTestTemplateCache() (map[string]*template.Template, error) {

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

