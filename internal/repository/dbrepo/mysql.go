package dbrepo

import (
	"context"
	"github.com/yaji1122/bookings-go/internal/model"
	"time"
)

func (m *mysqlDBRepo) AllUsers() bool {
	return true
}

func (m *mysqlDBRepo) InsertReservation(res model.Reservation) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	query := `insert into reservations (first_name, last_name, 
				email, phone, start_date, end_date, room_id, created_at, updated_at) 
			    values(?, ?, ?, ?, ?, ?, ?, ?, ?);`

	stmt, err := m.Pool.Prepare(query)
	checkErr(err)

	result, err := stmt.ExecContext(ctx, res.FirstName, res.LastName, res.Email, res.Phone, res.StartDate, res.EndDate, res.RoomID, time.Now(), time.Now())
	checkErr(err)

	id, err := result.LastInsertId()

	return id, nil
}

func (m *mysqlDBRepo) InsertRoomRestriction(res model.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, created_at, updated_at, restriction_id) 
				values (?, ?, ?, ?, ?, ?, ?)`

	_, err := m.Pool.ExecContext(ctx, query, res.StartDate, res.EndDate, res.RoomID, res.ReservationID, time.Now(), time.Now(), res.RestrictionID)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysqlDBRepo) SearchAvailabilityByDatesAndRoom(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var count int

	query := `SELECT count(id) FROM room_restriction WHERE room_id = ? AND ? < end_date AND ? > start_date`
	stmt, err := m.Pool.Prepare(query)
	checkErr(err)

	row := stmt.QueryRowContext(ctx, roomID, end, start)
	err = row.Scan(&count)
	checkErr(err)

	if count == 0 {
		return true, nil
	}
	return false, nil
}

func (m *mysqlDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]model.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var rooms []model.Room

	query := `SELECT r.id, r.room_name 
			  FROM rooms r 
			  WHERE r.id not in 
			        (SELECT room_id 
			         FROM room_restrictions rr WHERE ? < rr.end_date AND start > rr.start_date)`

	stmt, err := m.Pool.Prepare(query)
	checkErr(err)

	rows, err := stmt.QueryContext(ctx, end, start)
	checkErr(err)

	for rows.Next() {
		var room model.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		checkErr(err)
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
