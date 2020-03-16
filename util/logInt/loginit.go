package logInit

import (
	"log"
	"os"
	"time"
)

func init() {
	var dir = "log/"
	_, _ = EnsureFile(dir)
	filename := dir + time.Now().Format("20060102") + ".log"
	f, _ := os.OpenFile(filename,  os.O_RDWR| os.O_APPEND | os.O_CREATE, 0666)
	log.SetOutput(f)

	log.Println("程序已启动...")
}

func EnsureFile(file string) (*os.File, error){
	if !IsExist(file) {
		return  os.Create(file)
	}else{
		return os.OpenFile(file,  os.O_RDWR|os.O_CREATE, 0666)
	}
}

//IsExist  判断文件夹/文件是否存在  存在返回 true
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}