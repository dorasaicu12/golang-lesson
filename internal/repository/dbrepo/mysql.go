package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/dorasaicu12/booking/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *mysqlDBrepo) AllUsers() bool {
	return true
}
//insert reservation
func (m *mysqlDBrepo) InsertReservation(res models.Reservation)(int64,error) {
	ctx,cancle := context.WithTimeout(context.Background(),3*time.Second)
	defer cancle()
	stmt :=`insert into reservation (first_name,last_name,email,phone,start_date,end_date,room_id,created_at,updated_at)
	 values(?,?,?,?,?,?,?,?,?)
	`
    last,err:=	m.DB.ExecContext(ctx,stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)
	inserted,err:=last.LastInsertId()
	if(err != nil){
		return 0,err
	}
	return inserted,nil
}
func (m *mysqlDBrepo) InsertRoomRestriction(r models.RoomRestriction) error{
	ctx,cancle := context.WithTimeout(context.Background(),10*time.Second)
	defer cancle()
	stmt :=`insert into room_restrictions (reservation_id,restriction_id,room_id,start_date,end_date)
	values(?,?,?,?,?)
   `
	_,err:=m.DB.ExecContext(ctx,stmt,
		r.ReservationID,
		r.ResrictionID,
		r.RoomID,
		r.StartDate,
		r.EndDate,
	)
	if(err !=nil){
		return err
	}
	return nil
}
func (m *mysqlDBrepo) Search_Avai_Bydate_ID(RoomId,start,end any)(bool,error){
	ctx,cancle := context.WithTimeout(context.Background(),3*time.Second)
	defer cancle()
	var numrow int
	query :=`SELECT count(id) FROM room_restrictions WHERE room_id=? and start_date BETWEEN ? and ?`
	row :=m.DB.QueryRowContext(ctx,query,RoomId,start,end)
    err:=row.Scan(&numrow)
	if(err !=nil){
		return false,err
	}
	if numrow==0{
		return false,nil
	}
	return true,nil
}

func (m *mysqlDBrepo) Search_Avai_All_Room(start,end any)([]models.Room,error){
	ctx,cancle := context.WithTimeout(context.Background(),3*time.Second)
	defer cancle()
	var numrow []models.Room
	query :=`select r.id,r.room_name from  rooms as r where r.id  in 
	(SELECT rr.room_id from  room_restrictions as rr  where (rr.end_date BETWEEN ? and ?) and (rr.start_date BETWEEN ? and ?))`
	row,err :=m.DB.QueryContext(ctx,query,start,end,start,end)
 
	if(err !=nil){
		return numrow,err
	}
	for row.Next(){
		var room models.Room
		err:=row.Scan(
                  &room.ID,
				  &room.RoomName,
		)
		if(err !=nil){
			return numrow,err
		}
		numrow=append(numrow,room)
	}

	return numrow,nil
}
func (m *mysqlDBrepo) GetRoomByID (id int) (models.Room,error){
	ctx,cancle := context.WithTimeout(context.Background(),3*time.Second)
	defer cancle()
	var room models.Room
	query:=`select id,room_name from rooms where id =?`
	row :=m.DB.	QueryRowContext(ctx,query,id)
	err :=row.Scan(
		&room.ID,
		&room.RoomName,
	)
	if err !=nil{
		return room,err
	}
	return room,err
}
//retrun user by id
func (m *mysqlDBrepo) GetuserById(id int)(models.Users,error){
	ctx,cancle := context.WithTimeout(context.Background(),3*time.Second)
	defer cancle()
	query := `select id,first_name,last_name,email,password,access_level,created_at,updated_at from users where id =?`
	row :=m.DB.QueryRowContext(ctx,query,id)
	var u models.Users
	err :=row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Accesslevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u,err
	}
	return u,nil
} 
//update user
func (m *mysqlDBrepo) UpdateUser(u models.Users)(error){
	ctx,cancle := context.WithTimeout(context.Background(),3*time.Second)
	defer cancle()
	query :=`update users set first_name=?,last_name=?,email=?,access_level =?,updated_at=?`
	_,err :=m.DB.ExecContext(ctx,query,
        u.FirstName,
		u.LastName,
		u.Email,
		u.Accesslevel,
		u.UpdatedAt,
	)
	if err != nil {
      return err
	}
	return nil
}
//authenticate user
func (m *mysqlDBrepo) Authenticate(email,tesPassword string)(int,string,error){
	ctx,cancle := context.WithTimeout(context.Background(),3*time.Second)
	defer cancle()

	var id int
	var hashedPassword string

	row:=m.DB.QueryRowContext(ctx,"select id,password from users where email= ? ",email)
	err :=row.Scan(
		&id,
		&hashedPassword,
	)
	// log.Println(email)
	//https://go.dev/play/p/uKMMCzJWGsW generate password here
	if err !=nil{
		return id,"",err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(tesPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0,"",errors.New("incorect Passworld")
	}else if err != nil{
		return 0,"",err
	}
	return id,hashedPassword,nil
}
