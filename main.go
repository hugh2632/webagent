package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"webagent/pkg/crawler"
	"webagent/pkg/setting"
	"webagent/routers"
	_ "webagent/util/logInt"
)

func main() {
	{
		crawler.IsDubug = setting.RunMode == "debug"
		crawler.Capacity = setting.Capacity
		crawler.Timeout = setting.PAGE_Wait
		crawler.Read_Dealy = setting.READ_Delay
		crawler.FirstPage = "http://localhost:" + strconv.Itoa(setting.HTTPPort) + "/"
		if runtime.GOOS == "windows" {
			_ = crawler.Instance()
		}
	}

	router := routers.InitRouter()
	router.StaticFS("/resource", http.Dir("./backup"))
	router.StaticFS("/js", http.Dir("./web/js"))
	router.LoadHTMLGlob("web/page/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "搜索",
		})
	})
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}
