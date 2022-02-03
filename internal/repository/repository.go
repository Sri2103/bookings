package repository

import (
	"time"

	"github.com/Sri2103/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(r models.RoomRestriction) error

	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomId int) (bool, error)

	SearchAvaialabilityForAllRooms(start, end time.Time) ([]models.Room, error)

	GetRoomByID(id int) (models.Room, error)
}
