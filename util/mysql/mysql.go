package mysqlutil

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"runtime/debug"
)

func NewMysql(datasource string, f func(*sql.DB)) (err error) {
	if p := recover(); p != nil {
		str, ok := p.(string)
		if ok {
			err = errors.New(str)
			log.Println(str)
			fmt.Println(str)
		} else {
			err = errors.New("panic")
		}
		debug.PrintStack()
	}
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		panic("打开数据库失败," + err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	f(db)
	return err
}
