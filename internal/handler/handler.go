package handler

import (
	"encoding/json"
	"github.com/yaji1122/bookings-go/internal/config"
	"github.com/yaji1122/bookings-go/internal/model"
	"github.com/yaji1122/bookings-go/internal/render"
	"log"
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

//Booking page
func (m *Repository) Booking(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	session := m.AppConfig.Session
	session.Put(r.Context(), "remoteIP", remoteIP)
	session.Put(r.Context(), "title", "booking")
	render.TemplateRenderer(w, r, "booking", &model.TemplateData{})
}

//Contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"
	remoteIP := m.AppConfig.Session.GetString(r.Context(), "remoteIP")
	stringMap["remoteIP"] = remoteIP
	//send the data
	render.TemplateRenderer(w, r, "contact", &model.TemplateData{
		StringMap: stringMap,
	})
}

//Index page
func (m *Repository) Index(w http.ResponseWriter, r *http.Request) {
	render.TemplateRenderer(w, r, "index", &model.TemplateData{})
}

//Room page
func (m *Repository) Room(w http.ResponseWriter, r *http.Request) {
	render.TemplateRenderer(w, r, "room", &model.TemplateData{})
}

//Reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.TemplateRenderer(w, r, "reservation", &model.TemplateData{})
}

type jsonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

//Availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	//queryDate := r.Form.Get("queryDate")
	//stay := r.Form.Get("stay")

	resp := jsonResponse{
		Success: true,
		Message: "Available.",
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(out)

	//render.TemplateRenderer(w, r, "availability", &model.TemplateData{})
}

//func About(w http.ResponseWriter, r *http.Request) {
//	sum, _ := addValues(2, 3)
//	_, _ = fmt.Fprintf(w, fmt.Sprintf("This is the about page and 2 + 3 is %d", sum))
//}
