package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"

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
    defer close(app.MailChan)
	listenForMail()
	fmt.Println(fmt.Sprintf("Listen for mail function at port :8025"))
	// msg :=models.MailData{
	// 	To: "JONHN@gmail.com",
	// 	From: "me@gmail.com",
	// 	Subject: "Somw subject",
	// 	Content: "",
	// }
	// app.MailChan <- msg
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

		mailChan :=make(chan models.MailData)
		app.MailChan=mailChan
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
		db,err :=driver.ConnectSQL("root:@(127.0.0.1:3306)/booking?parseTime=true")
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

	mux.With(Auth).Get("/", handlers.Repo.Home)
	mux.With(Auth).Route("/", func(mux chi.Router) {
		mux.Get("/about", handlers.Repo.About)
		mux.Get("/general-quate", handlers.Repo.General)
		mux.Get("/major-suite", handlers.Repo.Major)
	
		mux.Get("/search-avai", handlers.Repo.Search_Avai)
		mux.Post("/search-avai", handlers.Repo.Handle_Search_Avai)
		mux.Post("/search-avai-json", handlers.Repo.AvaibilityJSON)
		mux.Get("/choose-room/{id}",handlers.Repo.ChooseRoom)
		mux.Get("/book-room",handlers.Repo.BookRoom)
	})
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summery",handlers.Repo.ReservationSummery)

	mux.Get("/user-login",handlers.Repo.LoginPage)
	mux.Post("/user-login",handlers.Repo.PostLogin)

	mux.Get("/user-logout",handlers.Repo.Logout)

	mux.Route("/admin",func(mux chi.Router){
		//mux.Use(Auth)
		mux.Get("/dashboard",handlers.Repo.AdminDashboard)
		mux.Get("/new-reservation",handlers.Repo.AdminNewReservation)
		mux.Get("/all-reservation",handlers.Repo.AdminAllReservation)
		mux.Get("/reservation-calendar",handlers.Repo.AdminCalendarReservation)

		mux.Get("/reservation-show/{src}/{id}",handlers.Repo.AdminShowRevervation)
		mux.Post("/reservation-show/{src}/{id}",handlers.Repo.AdminPostShowRevervation)

		mux.Get("/processed/{src}/{id}",handlers.Repo.AdminProcessedReservation)
		mux.Get("/deleted/{src}/{id}",handlers.Repo.AdminDeletedReservation)
	})

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
func Auth(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if !helpers.IsAuthenticate(r){
			session.Put(r.Context(),"error","Please Login first")
			http.Redirect(w,r,"/user-login",http.StatusSeeOther)
		}
		next.ServeHTTP(w,r)
	})
}

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMessage(msg)
		}
	}()
}

func sendMessage(m models.MailData){
	server :=mail.NewSMTPClient()
	server.Host="localhost"
	server.Port=1025
	server.KeepAlive=false
	server.ConnectTimeout=10*time.Second
	server.SendTimeout=10*time.Second

	client,err :=server.Connect()
	if err !=nil {
		errorLog.Println(err)
	}
	email:=mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == ""{
		email.SetBody(mail.TextHTML,string(m.Content))
	}else{
		data,err := ioutil.ReadFile(fmt.Sprintf("./email-template/%s",m.Template))
		if err !=nil{
			app.ErrorLog.Println(err)
		}
		mailtemplate :=string(data)
		msgToSend := strings.Replace(mailtemplate,"[%body%]",string(m.Content),1)
		email.SetBody(mail.TextHTML,msgToSend)
	}
	// email.SetBody(mail.TextHTML,string(m.Content))
	err =email.Send(client)
	if err !=nil{
		log.Println(err)
	}else{
		log.Println("email send")
	}
}
