package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFrom_Valid(t *testing.T) {
	request := httptest.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	request := httptest.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required filed missing")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "a")
	postData.Add("c", "a")

	r, _ := http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows doesn't have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	request := httptest.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)

	if form.Has("a") {
		t.Error("form shows has value when required filed missing")
	}

	postData := url.Values{}
	postData.Add("a", "")

	r := httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postData
	form = New(r.PostForm)

	if form.Has("a") {
		t.Error("form shows has value when required filed is blank")
	}

	postData2 := url.Values{}
	postData2.Add("a", "aaaa")
	g := httptest.NewRequest("POST", "/whatever", nil)
	g.PostForm = postData2
	form = New(g.PostForm)

	if !form.Has("a") {
		t.Error("form shows no value when required filed is filled")
	}

}
