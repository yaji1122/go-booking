package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	Key   string
	Value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"booking", "/booking", "GET", []postData{}, http.StatusOK},
	{"room", "/room", "GET", []postData{}, http.StatusOK},
	{"reservation", "/reservation", "GET", []postData{}, http.StatusOK},
	{"reservation", "/reservation", "POST", []postData{
		{Key: "firstName", Value: "Yachi"},
		{Key: "lastName", Value: "Lin"},
		{Key: "email", Value: "test@gmail.com"},
		{Key: "bookingDate", Value: "2021/09/28"},
		{Key: "phone", Value: "0912345678"},
		{Key: "stay", Value: "2"},
	}, http.StatusOK},
	{"summary", "/summary", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "POST", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			response, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		} else {
			values := url.Values{}

			for _, x := range e.params {
				values.Add(x.Key, x.Value)
			}

			response, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
