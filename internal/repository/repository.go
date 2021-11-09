package repository

import "github.com/yaji1122/bookings-go/internal/model"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res model.Reservation) error
}
