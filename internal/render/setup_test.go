package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dorasaicu12/booking/internal/config"
	"github.com/dorasaicu12/booking/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M){
			// what i am going to put the session
			gob.Register(models.Reservation{})
			// change this to true when in production
			testApp.InProduction = false
		
			// set up the session
			session = scs.New()
			session.Lifetime = 24 * time.Hour
			session.Cookie.Persist = true
			session.Cookie.SameSite = http.SameSiteLaxMode
			session.Cookie.Secure = false
		
			testApp.Session = session
			app=&testApp
	os.Exit(m.Run())
}
type myWritter struct{
    
}
func (tw *myWritter) Header() http.Header{
	var h http.Header
	return h
}
func (tw *myWritter) WriteHeader (i int){

}
func (tw *myWritter) Write(b []byte)(int,error){
	lenght:=len(b)
	return lenght,nil
}