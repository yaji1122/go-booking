package repository

import (
	"github.com/yaji1122/bookings-go/internal/model"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res model.Reservation) (int64, error)
	InsertRoomRestriction(res model.RoomRestriction) error
	SearchAvailabilityByDatesAndRoom(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]model.Room, error)
}
