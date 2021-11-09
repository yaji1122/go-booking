package render

import (
	"github.com/yaji1122/bookings-go/internal/model"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td model.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	result := AddDefaultData(&td, r)

	if result == nil {
		t.Error("failed")
	}
}

func TestTemplateRenderer(t *testing.T) {
	rootPath = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	appConfig.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var w myWriter

	err = Template(&w, r, "index", &model.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(&w, r, "contact", &model.TemplateData{})
	if err != nil {
		t.Error("rendered template that doesnt exist.")
	}
}

func TestCreateTemplateCache(t *testing.T) {
	rootPath = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}
