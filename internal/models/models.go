package models

import (
	"html/template"
	"time"
)


type Users struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	PassWord    string
	Accesslevel int
	CreatedAt   time.Time
	UpdatedAt time.Time
}
type Room struct {
	ID int
	RoomName string
	CreatedAt   time.Time
	UpdatedAt time.Time
}
type Restriction struct {
	ID int
	RestrictionName string
	CreatedAt   time.Time
	UpdatedAt time.Time
}
type Reservation struct{
	ID int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate   time.Time
	EndDate time.Time
	RoomID int
	CreatedAt   time.Time
	UpdatedAt time.Time
	Room Room
}
type RoomRestriction struct{
	ID int
	StartDate   time.Time
	EndDate time.Time
	CreatedAt   time.Time
	UpdatedAt time.Time
	RoomID int
	ReservationID int64
	ResrictionID int
	Room Room
	Reservation Reservation
}

//mail data
type MailData struct {
	To string
	From string
	Subject string
	Content template.HTML
	Template string
}