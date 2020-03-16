package mysqlutil

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func NewMysql(datasource string, f func(*sql.DB)) *error {
	defer func() {
		er := recover()
		if er != nil {
			log.Println(er.(error).Error())
		}
	}()
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	f(db)
	return &err
}





