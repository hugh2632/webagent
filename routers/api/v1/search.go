package v1

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"webagent/model"
	"webagent/pkg/setting"
	mysqlutil "webagent/util/mysql"
)

func Search(c *gin.Context) {
	w, _ := c.GetPostForm("key")
	//从es加载
	mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
		var datas []model.WebData
		rows, err := db.Query(`SELECT Taskres_pageurl, Taskres_pagetitle, Taskres_pagedate, Taskres_pagepath FROM taskres where Taskres_pagetitle like '%` + w + `%' or Taskres_pagetext like '%` + w + `%' and Taskres_status = 1`)
		if err != nil{
			log.Println(err)
		}
		for rows.Next(){
			var tmp model.WebData
			err = rows.Scan(&tmp.Webdata_url, &tmp.Webdata_title,&tmp.Webdata_date, &tmp.Webdata_path)
			datas = append(datas, tmp)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"iserror": true,
				"msg":     "搜索失败!" ,
			})
		} else {
			res, _ := json.Marshal(datas)
			c.JSON(http.StatusOK, gin.H{
				"iserror": false,
				"msg":     "搜索成功",
				"data":   string(res),
			})
		}
	})

}