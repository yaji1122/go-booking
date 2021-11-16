package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/alexedwards/scs/v2"
	"github.com/yaji1122/bookings-go/internal/forms"
	"github.com/yaji1122/bookings-go/internal/helper"
	"github.com/yaji1122/bookings-go/internal/logger"
	"github.com/yaji1122/bookings-go/internal/model"
	"github.com/yaji1122/bookings-go/internal/pageRenderer"
	"github.com/yaji1122/bookings-go/internal/repository"
	"github.com/yaji1122/bookings-go/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"time"
)

var log *logger.Logger
var session *scs.SessionManager
var dbRepository repository.DatabaseRepo

// CreateHandler sets the variables for the handlers
func CreateHandler(logger *logger.Logger, scs *scs.SessionManager, pool *sql.DB) {
	log = logger
	session = scs
	dbRepository = dbrepo.NewMysqlRepo(pool)
}

//Contact page
func Contact(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"
	remoteIP := session.GetString(r.Context(), "remoteIP")
	stringMap["remoteIP"] = remoteIP
	//send the data
	pageRenderer.Template(w, r, "contact", &model.TemplateData{
		StringMap: stringMap,
	})
}

//Index page
func Index(w http.ResponseWriter, r *http.Request) {
	pageRenderer.Template(w, r, "index", &model.TemplateData{})
}

//Room page
func Room(w http.ResponseWriter, r *http.Request) {
	pageRenderer.Template(w, r, "room", &model.TemplateData{})
}

//Reservation page
func Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation model.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	pageRenderer.Template(w, r, "reservation", &model.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func PostReservation(w http.ResponseWriter, r *http.Request) {
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

	//insert to db
	reservationID, err := dbRepository.InsertReservation(reservation)
	log.InfoLogger.Printf("ReservationID : %d", reservationID)
	restriction := model.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: reservationID,
		RestrictionID: 1,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}

	err = dbRepository.InsertRoomRestriction(restriction)

	if err != nil {
		helper.ServerError(w, err)
	}

	session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/summary", http.StatusSeeOther)

}

//Summary page
func Summary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := session.Get(r.Context(), "reservation").(model.Reservation)
	if !ok {
		log.ErrorLogger.Println("Cant get session.")
		session.Put(r.Context(), "error", "Can't get reservation from session.")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		log.InfoLogger.Println("cannot get item from session")
		return
	}

	session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	pageRenderer.Template(w, r, "summary", &model.TemplateData{Data: data})
}

type jsonResponse struct {
	Success bool         `json:"success"`
	Rooms   []model.Room `json:"rooms"`
}

//ajaxCheckAvailability get available rooms
func ajaxCheckAvailability(w http.ResponseWriter, r *http.Request) {
	shortForm := "1989/11/22"
	startDate, err := time.Parse(shortForm, r.Form.Get("startDate"))
	checkErr(err)
	endDate, err := time.Parse(shortForm, r.Form.Get("endDate"))

	rooms, err := dbRepository.SearchAvailabilityForAllRooms(startDate, endDate)
	checkErr(err)

	resp := jsonResponse{
		Success: true,
		Rooms:   rooms,
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helper.ServerError(w, err)
		return
	}
	w.Header().Set("Content-type", "application/json")
	_, err = w.Write(out)
	checkErr(err)

}

//Availability Page
func Availability(w http.ResponseWriter, r *http.Request) {
	pageRenderer.Template(w, r, "contact", &model.TemplateData{})
}

func checkErr(err error) {
	if err != nil {

	}
}
