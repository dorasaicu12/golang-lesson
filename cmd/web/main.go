package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dorasaicu12/booking/internal/config"
	"github.com/dorasaicu12/booking/internal/driver"
	"github.com/dorasaicu12/booking/internal/handlers"
	"github.com/dorasaicu12/booking/internal/helpers"
	"github.com/dorasaicu12/booking/internal/models"
	"github.com/dorasaicu12/booking/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
// main is the main function
func main() {
   db, err :=run()
    if err !=nil{
		log.Fatal(err)
	}
	defer db.SQL.Close()
	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB,error){
		// what i am going to put the session
		gob.Register(models.Reservation{})
		gob.Register(models.Users{})
		gob.Register(models.Room{})
		gob.Register(models.RoomRestriction{})
		// change this to true when in production
		app.InProduction = false
	    infoLog =log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime)
		app.InfoLog=infoLog
		errorLog=log.New(os.Stdout,"ERROR\t",log.Ldate|log.Ltime|log.Lshortfile)
		app.ErrorLog=errorLog
		// set up the session
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = app.InProduction
	
		app.Session = session

		//connect to database
        log.Println("connecting to database....")
		db,err :=driver.ConnectSQL("root:@(127.0.0.1:3306)/booking")
		if(err !=nil){
			log.Fatal(err)
			return nil,err
		}
		tc, err := render.CreateTemplateCache()
		if err != nil {
			log.Fatal("cannot create template cache")
			return nil,err
		}
	
		app.TemplateCache = tc
		app.UseCache = false
	
		repo := handlers.NewRepo(&app,db)
		handlers.NewHandlers(repo)
	    helpers.NewHelpers(&app)
		render.NewRendered(&app)
	return db,nil
}

// routing
func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/general-quate", handlers.Repo.General)
	mux.Get("/major-suite", handlers.Repo.Major)

	mux.Get("/search-avai", handlers.Repo.Search_Avai)
	mux.Post("/search-avai", handlers.Repo.Handle_Search_Avai)
	mux.Post("/search-avai-json", handlers.Repo.AvaibilityJSON)
	
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summery",handlers.Repo.ReservationSummery)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

//middleware

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
