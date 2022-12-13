package form

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct,embed a url value
type Form struct {
	url url.Values
	Errors errors
}

// new innitialize a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// VALID RETURN TRUE IF THERE IS NO ERROR
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
func (f *Form)Required(fields ...string){
   for _,field :=range fields{
	value :=f.url.Get(field)
	if strings.TrimSpace(value) ==""{
		f.Errors.Add(field, "This field can not be blank")
	}
   }
}

func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field can not be blank")
		return false
	}
	return true
}
func (f *Form) MinLenght(field string,lenght int,r * http.Request) bool{
	x := f.url.Get(field)
	if len(x)<lenght{
		f.Errors.Add(field,fmt.Sprintf("This character is at least %d character long",lenght))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string){
	if! govalidator.IsEmail(f.url.Get(field)){
		f.Errors.Add(field, "This is not email properly")
	}
}