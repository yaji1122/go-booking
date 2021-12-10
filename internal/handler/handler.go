package handler

import (
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/yaji1122/bookings-go/internal/helper"
	"github.com/yaji1122/bookings-go/internal/logger"
	"github.com/yaji1122/bookings-go/internal/mail"
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
var mailConfig *mail.Config

// CreateHandler sets the variables for the handlers
func CreateHandler(logger *logger.Logger, scs *scs.SessionManager, pool *sql.DB, config *mail.Config) {
	log = logger
	session = scs
	dbRepository = dbrepo.NewMysqlRepo(pool)
	mailConfig = config
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
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	checkErr(err)

	data := make(map[string]interface{})
	data["room"] = dbRepository.GetRoomById(id)

	pageRenderer.Template(w, r, "room", &model.TemplateData{Data: data})
}

//Reservation page
func Reservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helper.ServerError(w, err)
		return
	}

	// 01/02 03:04:05PM '06 -0700
	layout := "2006/01/02" //format of the date string
	startDate, err := time.Parse(layout, r.Form.Get("start_date"))
	checkErr(err)
	endDate, err := time.Parse(layout, r.Form.Get("end_date"))
	checkErr(err)
	roomID, err := strconv.Atoi(r.Form.Get("ID"))
	checkErr(err)
	room := dbRepository.GetRoomById(roomID)

	reservation := model.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation
	data["room"] = room

	pageRenderer.Template(w, r, "reservation", &model.TemplateData{
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
	layout := "2006/01/02" //format of the date string
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

	htmlMessage := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear %s: <br>
		This is confirm your reservation from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	//send notification
	msg := mail.Data{
		From:     "test@hog.com",
		To:       reservation.Email,
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	mailConfig.MailChan <- msg

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

//Availability Page
func Availability(w http.ResponseWriter, r *http.Request) {
	pageRenderer.Template(w, r, "availability", &model.TemplateData{})
}

type jsonResponse struct {
	Success bool         `json:"success"`
	Rooms   []model.Room `json:"rooms"`
}

//PostAvailability get available rooms
func PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helper.ServerError(w, err)
		return
	}

	layout := "2006/01/02"
	startDate, err := time.Parse(layout, r.Form.Get("start_date"))
	checkErr(err)
	endDate, err := time.Parse(layout, r.Form.Get("end_date"))
	checkErr(err)
	rooms, err := dbRepository.SearchAvailabilityForAllRooms(startDate, endDate)
	checkErr(err)

	if len(rooms) == 0 {
		session.Put(r.Context(), "msg", "No Vacancy.")
		http.Redirect(w, r, "/availability", http.StatusSeeOther)
		return
	} else {
		data := make(map[string]interface{})
		reservation := model.Reservation{
			StartDate: startDate,
			EndDate:   endDate,
		}
		data["rooms"] = rooms
		data["reservation"] = reservation
		pageRenderer.Template(w, r, "choose-room", &model.TemplateData{Data: data})
	}

	//out, err := json.MarshalIndent(resp, "", "    ")
	//if err != nil {
	//	helper.ServerError(w, err)
	//	return
	//}
	//w.Header().Set("Content-type", "application/json")
	//_, err = w.Write(out)
	//checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
