package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dorasaicu12/booking/internal/config"
	"github.com/dorasaicu12/booking/internal/driver"
	"github.com/dorasaicu12/booking/internal/form"
	"github.com/dorasaicu12/booking/internal/helpers"
	"github.com/dorasaicu12/booking/internal/models"
	"github.com/dorasaicu12/booking/internal/render"
	"github.com/dorasaicu12/booking/internal/repository"
	"github.com/dorasaicu12/booking/internal/repository/dbrepo"
	"github.com/go-chi/chi"
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
	layout :="1/2/2006"
	startDate,err :=time.Parse(layout,start)
	endDate,err :=time.Parse(layout,end)
	if err != nil{
		log.Fatal("ERROR IS:",err)
		return
	}
	rooms,err :=m.DB.Search_Avai_All_Room(startDate.Format("2006-01-02"),endDate.Format("2006-01-02"))
	if err !=nil{
		helpers.ServerError(w,err)
	}

     if len(rooms)==0{
		m.App.Session.Put(r.Context(),"error","No Room is AVailible")
		http.Redirect(w,r,"/search-avai",http.StatusSeeOther)
	 }else{
		data :=make(map[string]interface{})
		data["rooms"]=rooms
		res:=models.Reservation{
			StartDate: startDate,
			EndDate: endDate,
		}
		m.App.Session.Put(r.Context(),"reservation",res)
		render.Template(w, r,"chooseRoom.page.tmpl", &models.TemplateData{
		  Data: data,
	  })
	  return
	 }
     w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s",startDate.Format("2006-01-02"),endDate.Format("2006-01-02"))))
}
//handle send json respone
type jsonRespone struct {
	OK bool `json:"ok"`
	Message string `json:"message"`
	Room_id int `json:"room_id"`
	Start_date string `json:"start_date"`
	End_date string `json:"end_date"`
}
func (m *Repository) AvaibilityJSON(w http.ResponseWriter, r *http.Request) {

	sd :=r.Form.Get("start")
	ed :=r.Form.Get("end")
	layout :="1/2/2006"
	startDate,_ :=time.Parse(layout,sd)
	endDate,_ :=time.Parse(layout,ed)

	roomID,_:=strconv.Atoi(r.Form.Get("room_id"))
	available,err:=m.DB.Search_Avai_Bydate_ID(roomID,startDate,endDate)
	if err !=nil{
		helpers.ServerError(w,err)
	}
	resp := jsonRespone{
		OK:available,
		Message: "Available",
		Room_id:roomID,
		Start_date:sd,
		End_date:ed,
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
	res,foo :=m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !foo{
		helpers.ServerError(w,errors.New("Can not get reservation from session"))
		return

	}
	room,err :=m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	m.App.Session.Put(r.Context(),"reservation",res)
	res.Room.RoomName=room.RoomName
	sd :=res.StartDate.Format("2006-01-02")
	ed :=res.EndDate.Format("2006-01-02")
	data :=make(map[string]interface{})
	stringMap :=make(map[string]string)
	stringMap["start_date"]=sd
	stringMap["end_date"]=ed
	data["reservation"]= res
	render.Template(w, r,"make-reservation.page.tmpl", &models.TemplateData{
		Form:form.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}
//post reservation handle
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation,ok :=m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if ! ok{
      helpers.ServerError(w,errors.New("can not get the Session form post reservation"))
	  return
	}
     err :=r.ParseForm()
    
	 if err != nil{
		helpers.ServerError(w,err)
		return
	 }
	//  sd :=r.Form.Get("start_date")
	//  if err != nil{
	// 	helpers.ServerError(w,err)
	// 	return
	//  }
	//  ed :=r.Form.Get("end_date")
	//  if err != nil{
	// 	helpers.ServerError(w,err)
	// 	return
	//  }
	//  roomId,err := strconv.Atoi((r.Form.Get("room_id")))
	 //01/02 03:04:05PM '06 -0700
	//  layout :="1/2/2006 15-04-05"
	//  startDate,err :=time.Parse(layout,sd+" 00-00-00")
	//  endDate,err :=time.Parse(layout,ed+" 00-00-00")
	room,err :=m.DB.GetRoomByID(reservation.RoomID)
	 reservation.FirstName=r.Form.Get("first_name")
	 reservation.LastName=r.Form.Get("last_name")
	 reservation.Phone=r.Form.Get("phone")
	 reservation.Email=r.Form.Get("email")
	 reservation.Room.RoomName=room.RoomName

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

	  id,err :=m.DB.InsertReservation(reservation)
	  if err != nil{
		helpers.ServerError(w,err)
		return
	 }
     restriction :=models.RoomRestriction{
		StartDate: reservation.StartDate,
		EndDate: reservation.EndDate,
		RoomID: reservation.RoomID,
		ReservationID: id,
		ResrictionID: 2,
	 }
	 err =m.DB.InsertRoomRestriction(restriction)
	 if err != nil{
		helpers.ServerError(w,err)
		return
	 }
   htmlMessage :=fmt.Sprintf(`
   <strong> Reservation confirmation <strong>
   Dear %s:,</br>
   this is to confirm your reservation from %s to %s
   `,reservation.FirstName,reservation.StartDate.Format("2006-01-02"),reservation.EndDate.Format("2006-01-02"))
	 	msg :=models.MailData{
		To: reservation.Email,
		From: "me@gmail.com",
		Subject: "Reservation ConFirmation",
		Content: template.HTML(htmlMessage),
		Template: "basic.html",
	}
	m.App.MailChan <- msg

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
	fmt.Println(reservation.Room.RoomName)
	m.App.Session.Remove(r.Context(),"reservation")
	data :=make(map[string]interface{})
	data["reservation"]=reservation
	sd :=reservation.StartDate.Format("2006-01-02")
	ed :=reservation.EndDate.Format("2006-01-02")
	stringMap :=make(map[string]string)
	stringMap["start_date"]=sd
	stringMap["end_date"]=ed
	render.Template(w, r,"reservation-summery.page.tmpl", &models.TemplateData{
			Data: data ,
			StringMap: stringMap,
	})
}
func (m *Repository) ChooseRoom(w http.ResponseWriter,r *http.Request){
	roomId,err:=strconv.Atoi(chi.URLParam(r,"id"))
	if err !=nil{
		helpers.ServerError(w,err)
		return
	}
	m.App.Session.Get(r.Context(),"reservation")
	res,foo :=m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !foo{
		helpers.ServerError(w,err)
		return

	}
	res.RoomID=roomId
	m.App.Session.Put(r.Context(),"reservation",res)
	http.Redirect(w,r,"/make-reservation",http.StatusSeeOther)
}
func (m *Repository) BookRoom(w http.ResponseWriter,r *http.Request){
 //id ,s,e
 ID,_:=strconv.Atoi(r.URL.Query().Get("id"))
 ed :=r.URL.Query().Get("s")
 sd :=r.URL.Query().Get("e")
 layout :="1/2/2006"
 startDate,_ :=time.Parse(layout,sd)
 endDate,_ :=time.Parse(layout,ed)
 var res models.Reservation
 res.RoomID =ID
 res.StartDate=startDate
 res.EndDate=endDate
 m.App.Session.Put(r.Context(),"reservation",res)
 http.Redirect(w,r,"/make-reservation",http.StatusSeeOther)
}

func (m *Repository) LoginPage(w http.ResponseWriter,r *http.Request){
	render.Template(w, r,"login.page.tmpl", &models.TemplateData{
   Form: form.New(nil),
})
   }

func  (m *Repository) PostLogin(w http.ResponseWriter,r *http.Request){
	_=m.App.Session.RenewToken(r.Context())

	err :=r.ParseForm()
	if err !=nil{
		log.Println(err)
	}
	var email string
	var password string
	email = r.Form.Get("email")
	password = r.Form.Get("password")
	form :=form.New(r.PostForm)
	form.Required("email","password")
	form.IsEmail("email")
	if !form.Valid(){
		render.Template(w,r,"login.page.tmpl",&models.TemplateData{
			Form: form,
		})
		return
	}
	// log.Println(email)
	id,_,err :=m.DB.Authenticate(email,password)
	if err !=nil{
		log.Println(err)
		m.App.Session.Put(r.Context(),"error","Invalid login credential")
		http.Redirect(w,r,"/user-login",http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(),"user_id",id)
	m.App.Session.Put(r.Context(),"flash","Login Success")
	http.Redirect(w,r,"/",http.StatusSeeOther)
}
//logout
func (m *Repository) Logout(w http.ResponseWriter,r *http.Request){
	m.App.Session.Remove(r.Context(),"user_id")
	http.Redirect(w,r,"/user-login",http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter,r *http.Request){
	render.Template(w,r,"admin-dashboard.page.tmpl",&models.TemplateData{})
}
func (m *Repository) AdminNewReservation(w http.ResponseWriter,r *http.Request){
	reservation,err :=m.DB.AllReservation()
	if err !=nil{
		helpers.ServerError(w,err)
		return
	}
	data :=make(map[string]interface{})
	data["reservation"] =reservation
	render.Template(w,r,"admin-new-reservation.page.tmpl",&models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminAllReservation(w http.ResponseWriter,r *http.Request){
	reservation,err :=m.DB.AllReservation()
	if err !=nil{
		helpers.ServerError(w,err)
		return
	}
	data :=make(map[string]interface{})
	data["reservation"] =reservation
	render.Template(w,r,"admin-all-reservation.page.tmpl",&models.TemplateData{
		Data: data,
	})
}
func (m *Repository) AdminShowRevervation(w http.ResponseWriter,r *http.Request){
	explode :=strings.Split(r.RequestURI, "/")
	str := explode[4]
    str = strings.Replace(str, "%7D", "", -1)
    num, err := strconv.Atoi(str)
    if err != nil {
        fmt.Println("Failed to convert string to int")
        return
    }
	src :=explode[3]
	stringMap :=make(map[string]string)
	stringMap["src"]=src
	//get revervation
	res,err:=m.DB.GetOneReservation(num)
	data:=make(map[string]interface{})
	data["reservation"]=res
	render.Template(w,r,"admin-show.page.tmpl",&models.TemplateData{
		StringMap:stringMap,
		Data: data,
		Form: form.New(nil),
	})
}
func (m *Repository) AdminPostShowRevervation(w http.ResponseWriter,r *http.Request){
	err :=r.ParseForm()
	if err != nil{
	   helpers.ServerError(w,err)
	   return
	}
	explode :=strings.Split(r.RequestURI, "/")
	str := explode[4]
    str = strings.Replace(str, "%7D", "", -1)
	src :=explode[3]
    num, err := strconv.Atoi(str)
	if err != nil {
        fmt.Println("Failed to convert string to int")
        return
    }
	res,err := m.DB.GetOneReservation(num)
	if err != nil {
        helpers.ServerError(w,err)
        return
    }

	res.FirstName=r.Form.Get("first_name")
	res.LastName=r.Form.Get("last_name")
	res.Phone=r.Form.Get("phone")
	res.Email=r.Form.Get("email")
	err = m.DB.UpdateReservation(res,num)
	if err != nil {
        helpers.ServerError(w,err)
        return
    }
	m.App.Session.Put(r.Context(),"flash","Change saved")
	http.Redirect(w,r,fmt.Sprintf("/admin/%s-reservation",src),http.StatusSeeOther)
}

func (m *Repository) AdminCalendarReservation(w http.ResponseWriter,r *http.Request){
	render.Template(w,r,"admin-calendar.page.tmpl",&models.TemplateData{})
}

func (m *Repository) AdminProcessedReservation(w http.ResponseWriter,r *http.Request){
	id,_:=strconv.Atoi(chi.URLParam(r,"id"))
	src:=chi.URLParam(r,"src")
	err :=m.DB.UpdateProcessReservation(id,1)
	if err != nil {
        helpers.ServerError(w,err)
        return
    }
	m.App.Session.Put(r.Context(),"flash","Reservation mark ass resource")
	http.Redirect(w,r,fmt.Sprintf("/admin/%s-reservation",src),http.StatusSeeOther)
}

func (m *Repository) AdminDeletedReservation(w http.ResponseWriter,r *http.Request){
	id,_:=strconv.Atoi(chi.URLParam(r,"id"))
	src:=chi.URLParam(r,"src")
	err :=m.DB.DeleteReservation(id)
	if err != nil {
        helpers.ServerError(w,err)
        return
    }
	m.App.Session.Put(r.Context(),"flash","Reservation has ben deleted")
	http.Redirect(w,r,fmt.Sprintf("/admin/%s-reservation",src),http.StatusSeeOther)
}



