package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"general", "/general-quate", "GET", []postData{}, http.StatusOK},
	{"major-suite", "/major-suite", "GET", []postData{}, http.StatusOK},
	{"search-avai", "/search-avai", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"search-reservation", "/search-avai", "POST", []postData{
		 {key:"start",value: "2022-8-10"},
		 {key:"end",value: "2022-12-10"},
	}, http.StatusOK},
	{"search-reservation-json", "/search-avai-json", "POST", []postData{
		{key:"start",value: "2022-8-10"},
		{key:"end",value: "2022-12-10"},
   }, http.StatusOK},
   {"make-reservation", "/make-reservation", "POST", []postData{
	{key:"first_name",value: "che"},
	{key:"last_name",value: "anh"},
	{key:"email",value: "anh@gmail.com"},
	{key:"phone",value: "123123123"},
}, http.StatusOK},
}
func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
            values :=url.Values{}
			for _,x :=range e.params{
				values.Add(x.key,x.value)
			}
			resp,err := ts.Client().PostForm(ts.URL +e.url,values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
