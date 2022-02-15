package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/Sri2103/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

//InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var newID int
	stmt := `insert into reservations(first_name,last_name,email,phone,start_date,end_date,
		room_id,created_at,updated_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id `

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

//InsertRoomRestriction insert room restriction into a database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	stmt := `insert into room_restrictions (start_date,end_date,room_id,
		reservation_id,restriction_id,created_at,updated_at) 
		values
		($1,$2,$3,$4,$5,$6,$7) `

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

//SearchAvailabilityByDatesByRoomID return true if there is availability, returns false if there is no availability
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(roomId int, start, end time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var numRows int
	defer cancel()
	query := `
		select 
			count(id)
		from
			room_restrictions 
		where 
				room_id = $1
				and   $2<end_date and $3 > start_date;`

	row := m.DB.QueryRowContext(ctx, query, roomId, start, end)

	err := row.Scan(&numRows)

	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

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
func (m *postgresDBRepo) SearchAvaialabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		select r.id, r.room_name 
		from
			rooms r
		where r.id  not in(select room_id from room_restrictions rr where $1<rr.end_date and $2>rr.start_date);	
	`
	rows, err := m.DB.QueryContext(ctx, query, start, end)

	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		rows.Scan(
			&room.ID,
			&room.RoomName)
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}

//GetRoomByID gets a room by referring to the ID
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var room models.Room

	query := ` 
		select 
			id,room_name,created_at,updated_at 
		from 
			rooms 
		where
			 id = $1 ;`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}

	return room, nil

}

//GetUserByID returns user by id
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `select  id,first_name,last_name,email,password,access_level,created_at,updated_at 
	from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil

}

//UpdateUser updates the user values in database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `update users set first_name =$1,last_name =$2,email = $3,access_level = $4,created_at = $5,updated_at = $6 `

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil

}

//Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id,password from users where email=$1", email)

	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("Incorrect Password")
	} else if err != nil {
		return 0, "", err

	}

	return id, hashedPassword, nil

}

//AllReservations returns a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id,r.first_name,r.last_name,r.email,r.phone,r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at,r.processed,rm.id,rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		order by r.start_date asc`
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err

	}

	return reservations, nil

}

//AllNewReservations reservations that are not processed
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		select r.id,r.first_name,r.last_name,r.email,r.phone,r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at,rm.id,rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id)
		where processed=0
		order by r.start_date asc`
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)

		if err != nil {
			return reservations, err
		}

		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err

	}

	return reservations, nil

}

//GetReservationById returns One reservation   by ID
func (m *postgresDBRepo) GetReservationById(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var res models.Reservation

	query := `
	select r.id,r.first_name,r.last_name,r.email,r.phone,r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at,r.processed,rm.id,rm.room_name
	from reservations r
	left join rooms rm on (r.room_id = rm.id)
	where r.id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)

	if err != nil {
		return res, err
	}

	return res, nil

}

//UpdateReservation updates the reservation  in database
func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `update reservations set first_name =$1,last_name =$2,email = $3,phone= $4,updated_at = $5 wher id = $6`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}

	return nil

}

//DeleteReservation deletes a reservation by id  from database
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := ` delete from reservations where id = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil

}

//UpdateProcessed for reservation for id
func (m *postgresDBRepo) UpdateProcessed(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `update reservations set processed = $1 where id =$2`

	_, err := m.DB.ExecContext(ctx, query, id, processed)
	if err != nil {
		return err
	}

	return nil

}

func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `select id,room_name,created_at,updated_at from rooms order by room_name`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return rooms, err
	}

	defer rows.Close()

	for rows.Next() {
		var rm models.Room
		err = rows.Scan(
			&rm.ID,
			&rm.RoomName,
			&rm.CreatedAt,
			&rm.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, rm)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}

//GetRestrictionsForRoomByDate returns restrictions for a room by date range
func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomId int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction
	query := `select id, coalesce(reservation_id,0),restriction_id,room_id,start_date,end_date 
	from room_restrictions 
	where  $1<end_date and $2 >= start_date and room_id =$3`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomId)

	if err != nil {
		return restrictions, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err = rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return restrictions, err
		}
		restrictions = append(restrictions, r)
	}
	if err = rows.Err(); err != nil {
		return restrictions, err
	}
	return restrictions, nil
}
