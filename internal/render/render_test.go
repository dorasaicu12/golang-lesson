package render

import (
	"net/http"
	"testing"

	"github.com/dorasaicu12/booking/internal/models"
)

func TestAddDefault(t *testing.T){
	var td models.TemplateData
    r,err :=getSession()
	if err !=nil {
		t.Error(err)
	}
	session.Put(r.Context(),"flash","123")
	result := AddDefaultData(&td,r)
	if result.Flash !="123"{
		t.Error("flash value of 123 not found in session")
	}
}
func TestRenderTemplate(t *testing.T){
   pathToTemplate="./../../templates"
   tc,err :=CreateTemplateCache()
   if err != nil{
	t.Error(err)
   }
   r,err :=getSession()
   if err != nil{
	t.Error(err)
   }
  var ww myWritter
   app.TemplateCache=tc
   err =Template(&ww,r,"home.page.tmpl",&models.TemplateData{})
   if err != nil{
	  t.Error("error writting templte to browser")
   }
}
func getSession()(*http.Request,error){
	r,err :=http.NewRequest("GET","/some-url",nil)
	if err != nil{
		return nil,err
	}

	ctx :=r.Context()
	ctx,_=session.Load(ctx,r.Header.Get("X-Session"))
	r =r.WithContext(ctx)
	return r,nil
}