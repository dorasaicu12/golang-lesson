package repository

import "github.com/dorasaicu12/booking/internal/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int64,error)
	InsertRoomRestriction(r models.RoomRestriction) error
	Search_Avai_Bydate_ID(RoomId,start,end any)(bool,error)
	Search_Avai_All_Room(start,end any)([]models.Room,error)
	GetRoomByID(id int)(models.Room,error)
	GetuserById(id int)(models.Users,error)
	UpdateUser(u models.Users)(error)
	Authenticate(email,tesPassword string)(int,string,error)
	AllReservation()([]models.Reservation,error)
	GetOneReservation(id int)(models.Reservation,error)
	UpdateReservation(u models.Reservation,id int)(error)
	DeleteReservation(id int)(error)
	UpdateProcessReservation(id,processed int)(error)
}