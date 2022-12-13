package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dorasaicu12/booking/internal/config"
	"github.com/dorasaicu12/booking/internal/driver"
	"github.com/dorasaicu12/booking/internal/form"
	"github.com/dorasaicu12/booking/internal/helpers"
	"github.com/dorasaicu12/booking/internal/repository"
	"github.com/dorasaicu12/booking/internal/repository/dbrepo"

	"github.com/dorasaicu12/booking/internal/models"
	"github.com/dorasaicu12/booking/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig,db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB :dbrepo.NewMysqlRepo(db.SQL,a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	
	render.Template(w,r, "home.page.tmpl" ,&models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic


	// send data to the template
	render.Template(w, r,"about.page.tmpl", &models.TemplateData{
		
	})
}

func (m *Repository) General(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// send data to the template
	render.Template(w, r,"general.page.tmpl", &models.TemplateData{
	
	})
}
func (m *Repository) Major(w http.ResponseWriter, r *http.Request) {
	// perform some logic


	// send data to the template
	render.Template(w, r,"major.page.tmpl", &models.TemplateData{
		
	})
}

func (m *Repository) Search_Avai(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// send data to the template
	render.Template(w, r,"search-avai.page.tmpl", &models.TemplateData{
		
	})
}

func (m *Repository) Handle_Search_Avai(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

     w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s",start,end)))
}
//handle send json respone
type jsonRespone struct {
	OK bool `json:"ok"`
	Message string `json:"message"`
}
func (m *Repository) AvaibilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonRespone{
		OK:true,
		Message: "Available",
	}
	out,err := json.MarshalIndent(resp,"","   ")
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(out)
}


func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// send data to the template
	render.Template(w, r,"contact.page.tmpl", &models.TemplateData{
		
	})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// send data to the template
	var emptyReservation models.Reservation
	data :=make(map[string]interface{})
	data["reservation"]= emptyReservation
	render.Template(w, r,"make-reservation.page.tmpl", &models.TemplateData{
		Form:form.New(nil),
		Data: data,
	})
}
//post reservation handle
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
     err :=r.ParseForm()

	 if err != nil{
		helpers.ServerError(w,err)
		return
	 }
	 reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName: r.Form.Get("last_name"),
		Email: r.Form.Get("email"),
		Phone: r.Form.Get("phone"),
	 }
	 forms := form.New(r.PostForm)
	//forms.Has("first_name",r)
	forms.Required("first_name","last_name","email","phone")
    forms.MinLenght("first_name",5,r)
	forms.IsEmail("email")
	 if(!forms.Valid()){
          data :=make(map[string]interface{})
		  data["reservation"]=reservation
		  render.Template(w, r,"make-reservation.page.tmpl", &models.TemplateData{
			Form:forms,
			Data: data,
		})
		return
	 }
	 m.App.Session.Put(r.Context(),"reservation",reservation)
	 
	 //redirect
	 m.App.Session.Put(r.Context(),"flash","Create reservation successfully")
	 http.Redirect(w,r,"/reservation-summery",http.StatusSeeOther)
}

func (m *Repository) ReservationSummery(w http.ResponseWriter,r *http.Request){
	reservation,ok :=m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok{
		m.App.ErrorLog.Println("Can't get error from session")
       log.Println("can not get the session")
	   m.App.Session.Put(r.Context(),"error","can not get the session")
	   http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
	   return
	}
	data :=make(map[string]interface{})
	data["reservation"]=reservation
	render.Template(w, r,"reservation-summery.page.tmpl", &models.TemplateData{
			Data: data ,
	})
}
