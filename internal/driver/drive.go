package driver

import (
	"database/sql"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

//DB HOLD DATABASE CONNECTION
type DB struct {
	SQL *sql.DB
}
var dbCon =&DB{}
const maxOpenDbConn=10
const maxIdleDb=5
const maxDbLifetime=5 *time.Minute

//create databse pool for mysql
func ConnectSQL(dsn string)(*DB,error){
    d,err :=NewDataBase(dsn)
	if err !=nil{
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDb)
	d.SetConnMaxLifetime(maxDbLifetime)

	dbCon.SQL=d
	err =testDB(d)
	if err !=nil{
		return nil,err
	}
	return dbCon,nil
}
//try to ping the database
func testDB(d *sql.DB) error{
	err :=d.Ping()
	if err !=nil{
		return err
	}
	return nil
}
// Create a new databse connection
func NewDataBase(dsn string)(*sql.DB,error){
   db,err:=sql.Open("mysql",dsn)
   if err != nil{
	return nil,err
   }
   if err =db.Ping();err !=nil{
       return nil,err
   }
   return db,nil
}