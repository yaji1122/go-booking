package handler

import (
	"encoding/json"
	"github.com/yaji1122/bookings-go/internal/config"
	"github.com/yaji1122/bookings-go/internal/driver"
	"github.com/yaji1122/bookings-go/internal/forms"
	"github.com/yaji1122/bookings-go/internal/helper"
	"github.com/yaji1122/bookings-go/internal/model"
	"github.com/yaji1122/bookings-go/internal/render"
	"github.com/yaji1122/bookings-go/internal/repository"
	"github.com/yaji1122/bookings-go/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	AppConfig *config.AppConfig
	Pool      repository.DatabaseRepo
}

//NewRepo creates a new repository
func NewRepo(appConfig *config.AppConfig, pool *driver.Pool) *Repository {
	return &Repository{
		AppConfig: appConfig,
		Pool:      dbrepo.NewMysqlRepo(pool.SQL, appConfig),
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
	render.Template(w, r, "booking", &model.TemplateData{})
}

//Contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"
	remoteIP := m.AppConfig.Session.GetString(r.Context(), "remoteIP")
	stringMap["remoteIP"] = remoteIP
	//send the data
	render.Template(w, r, "contact", &model.TemplateData{
		StringMap: stringMap,
	})
}

//Index page
func (m *Repository) Index(w http.ResponseWriter, r *http.Request) {
	m.Pool.AllUsers()
	render.Template(w, r, "index", &model.TemplateData{})
}

//Room page
func (m *Repository) Room(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "room", &model.TemplateData{})
}

//Reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation model.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "reservation", &model.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helper.ServerError(w, err)
		return
	}

	// 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02" //format of the date string
	startDate, err := time.Parse(layout, r.Form.Get("start_date"))
	endDate, err := time.Parse(layout, r.Form.Get("end_date"))

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))

	reservation := model.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	//form.Has("firstName", r)
	form.Required("firstName", "lastName", "email", "bookingDate")

	if !reservation.IsValid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		err := render.Template(w, r, "reservation", &model.TemplateData{
			Data: data,
			Form: form,
		})
		if err != nil {
			return
		}
		return
	}
	//insert to db
	err = m.Pool.InsertReservation(reservation)
	if err != nil {
		helper.ServerError(w, err)
	}

	m.AppConfig.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/summary", http.StatusSeeOther)

}

//Summary page
func (m *Repository) Summary(w http.ResponseWriter, r *http.Request) {
	session := m.AppConfig.Session
	reservation, ok := session.Get(r.Context(), "reservation").(model.Reservation)
	if !ok {
		m.AppConfig.ErrorLog.Println("Cant get session.")
		session.Put(r.Context(), "error", "Can't get reservation from session.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		log.Println("cannot get item from session")
		return
	}

	session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "summary", &model.TemplateData{Data: data})
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
		helper.ServerError(w, err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(out)

	//render.Template(w, r, "availability", &model.TemplateData{})
}

//func About(w http.ResponseWriter, r *http.Request) {
//	sum, _ := addValues(2, 3)
//	_, _ = fmt.Fprintf(w, fmt.Sprintf("This is the about page and 2 + 3 is %d", sum))
//}
