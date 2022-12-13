package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/dorasaicu12/booking/internal/config"
	"github.com/go-chi/chi"
)

func TestMain(m *testing.M){
	os.Exit(m.Run())
}
type myHandler struct{
 handler http.Handler
}
func(mh *myHandler) serveHTTP(w http.ResponseWriter,r *http.Request){
	
}
//test main function
func TestRun(t *testing.T){
	_,err :=run()
	if err!=nil {
		t.Error("failed run")
	}
}
//test middleware
func TestNoSurf(t *testing.T){
	var myH myHandler
	h:=NoSurf(myH.handler)
	
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler but is %T", v))
	}
}

func TestSessionLoad(t *testing.T){
	var myH myHandler
	h:=SessionLoad(myH.handler)
	
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler but is %T", v))
	}
}

//etst route

func TestRoute(t *testing.T){
  var app config.AppConfig
  mux :=routes(&app)
  switch v:=mux.(type){
  case *chi.Mux:
	//do nothing
  default:
	t.Error(fmt.Sprintf("type is not *chi.mux but is %T", v))
  }
}

