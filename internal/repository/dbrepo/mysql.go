package dbrepo

import (
	"context"
	"github.com/yaji1122/bookings-go/internal/model"
	"time"
)

func (m *mysqlDBRepo) AllUsers() bool {
	return true
}

func (m *mysqlDBRepo) InsertReservation(res model.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	query := `insert into reservations (first_name, last_name, 
			email, phone, start_date, end_date, room_id, created_at, updated_at
			values (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := m.Pool.ExecContext(ctx, query, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return nil
}
