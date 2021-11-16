package pageRenderer

import (
	"github.com/yaji1122/bookings-go/internal/model"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td model.TemplateData
	rootPath = "./../../templates"
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
	configuration.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var w myWriter
	Template(&w, r, "index", &model.TemplateData{})
	Template(&w, r, "contact", &model.TemplateData{})
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
