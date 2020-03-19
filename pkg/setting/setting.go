package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort         int
	FirstPageTimeOut time.Duration
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration

	Capacity   int
	PAGE_Wait  time.Duration
	READ_Delay time.Duration

	MysqlDataSource string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadBase()
	LoadServer()
	LoadGYHLW()
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(30)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(120)) * time.Second
	PAGE_Wait = time.Duration(sec.Key("PAGE_Wait").MustInt(10)) * time.Second
	READ_Delay = time.Duration(sec.Key("READ_Delay").MustInt(0)) * time.Millisecond
}

func LoadGYHLW() {
	MysqlDataSource = Cfg.Section("mysql").Key("datasource").Value()
	Capacity = Cfg.Section("crawler").Key("capacity").MustInt(10)
}
