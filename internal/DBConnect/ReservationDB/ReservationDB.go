package ReservationDB

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
)
  var settings = mysql.ConnectionURL{
	Database: `booking`,
	Host:     `127.0.0.1:3306`,
	User:     `root`,
	Password: ``,
  }
  func DB() (db.Session,error){
	sess, err := mysql.Open(settings)
	if err !=nil {
		return nil,err
	}
	defer sess.Close()
	return sess,nil
  }