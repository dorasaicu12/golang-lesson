package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dorasaicu12/booking/internal/config"
	"github.com/dorasaicu12/booking/internal/models"
	"github.com/dorasaicu12/booking/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w,r, "home.page.tmpl" ,&models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send data to the template
	render.RenderTemplate(w, r,"about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) General(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// send data to the template
	render.RenderTemplate(w, r,"general.page.tmpl", &models.TemplateData{
	
	})
}
func (m *Repository) Major(w http.ResponseWriter, r *http.Request) {
	// perform some logic


	// send data to the template
	render.RenderTemplate(w, r,"major.page.tmpl", &models.TemplateData{
		
	})
}

func (m *Repository) Search_Avai(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// send data to the template
	render.RenderTemplate(w, r,"search-avai.page.tmpl", &models.TemplateData{
		
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
		log.Println(err)
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(out)
}


func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// send data to the template
	render.RenderTemplate(w, r,"contact.page.tmpl", &models.TemplateData{
		
	})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// send data to the template
	render.RenderTemplate(w, r,"make-reservation.page.tmpl", &models.TemplateData{
		
	})
}
