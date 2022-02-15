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
	if res.RoomID == 2 {
		return 0, errors.New("Some error")
	}

	return 1, nil
}

//InsertRoomRestriction insert room restriction into a database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
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

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil

}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 0, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	return reservations, nil
}

func (m *testDBRepo) GetReservationById(id int) (models.Reservation, error) {
	var res models.Reservation

	return res, nil
}

func (m *testDBRepo) UpdateReservation(u models.Reservation) error {
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {
	return nil
}

func (m *testDBRepo) UpdateProcessed(id, processed int) error {
	return nil
}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil

}

func (m *testDBRepo) GetRestrictionsForRoomByDate(roomId int, startDate, endDate time.Time) ([]models.RoomRestriction, error) {
	var restriction []models.RoomRestriction

	return restriction, nil
}
