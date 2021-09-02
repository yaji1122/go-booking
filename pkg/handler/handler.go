package handler

import (
	"github.com/yaji1122/bookings-go/pkg/config"
	"github.com/yaji1122/bookings-go/pkg/model"
	"github.com/yaji1122/bookings-go/pkg/render"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	AppConfig *config.AppConfig
}

//NewRepo creates a new repository
func NewRepo(appConfig *config.AppConfig) *Repository {
	return &Repository{
		AppConfig: appConfig,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.AppConfig.Session.Put(r.Context(), "remoteIP", remoteIP)
	render.TemplateRenderer(w, "home", &model.TemplateData{})
}

//About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"
	remoteIP := m.AppConfig.Session.GetString(r.Context(), "remoteIP")
	stringMap["remoteIP"] = remoteIP
	//send the data
	render.TemplateRenderer(w, "about", &model.TemplateData{
		StringMap: stringMap,
	})
}

//func About(w http.ResponseWriter, r *http.Request) {
//	sum, _ := addValues(2, 3)
//	_, _ = fmt.Fprintf(w, fmt.Sprintf("This is the about page and 2 + 3 is %d", sum))
//}
