package form

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r:=httptest.NewRequest("POST","/thaterve",nil)
    form:=New(r.PostForm)

	isValid :=form.Valid()
	if !isValid{
		t.Error("GOT A INVALID WHEN SHOULD HAVE BEEN")
	}
}

func TestForm_Required(t *testing.T) {
	r:=httptest.NewRequest("POST","/thaterve",nil)
    form:=New(r.PostForm)
    form.Required("a","b","c")
	if form.Valid(){
		t.Error("missing required fields")
	}
	postedData :=url.Values{}
	postedData.Add("a","a")
	postedData.Add("b","b")
	postedData.Add("c","c")
	r,_=http.NewRequest("POST","/thaterve",nil)
	r.PostForm=postedData
	form=New(r.PostForm)
	form.Required("a","b","c")
	if !form.Valid(){
		t.Error("show does not have required data")
	}
}


func TestForm_Min(t *testing.T) {
	r:=httptest.NewRequest("POST","/thaterve",nil)
    form:=New(r.PostForm)
    form.MinLenght("asdfsdfsdf",5,r)
	if form.Valid(){
		t.Error("missing required fields")
	}
	postedData :=url.Values{}
	postedData.Add("a","aq123123")
	r,_=http.NewRequest("POST","/thaterve",nil)
	r.PostForm=postedData
	form=New(r.PostForm)
	form.MinLenght("a",5,r)
	if !form.Valid(){
		t.Error("the min legnht is false")
	}
}
func TestForm_Email(t *testing.T){
   r :=httptest.NewRequest("POST","/WHATERVER",nil)
   form:=New(r.PostForm)
   form.IsEmail("a")
   if form.Valid(){
	t.Error("missing required fields")
   }
   postedData :=url.Values{}
   postedData.Add("a","asdasd@gmail.com")
   r,_ =http.NewRequest("POST","/WHATERVER",nil)
   r.PostForm=postedData
   form=New(r.PostForm)
   form.IsEmail("a")
   if !form.Valid(){
	t.Error("the field is not a email")
   }
}