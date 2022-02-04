package dbrepo

import (
	"errors"
	"time"

	"github.com/Sri2103/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

//InsertRoomRestriction insert room restriction into a database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	return nil
}

//SearchAvailabilityByDatesByRoomID return true if there is availability, returns false if there is no availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(roomId int, start, end time.Time) (bool, error) {

	return false, nil

}

// select
// 			count(id)
// 		from
// 			room_restrictions
// 		where
// 			room_id = $1
// 			and
// 			'$2' < end_date and '$3'> start_date;

//SearchAvaialabilityForAllRooms returns a slice of rooms if there are any.
func (m *testDBRepo) SearchAvaialabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil

}

//GetRoomByID gets a room by referring to the ID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id > 2 {
		return room, errors.New("room does not exist")
	}

	return room, nil

}
