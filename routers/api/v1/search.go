package v1

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"html"
	"log"
	"net/http"
	"webagent/model"
	"webagent/pkg/setting"
	mysqlutil "webagent/util/mysql"
)

func Search(c *gin.Context) {
	w, _ := c.GetPostForm("key")
	//从es加载
	_ = mysqlutil.NewMysql(setting.MysqlDataSource, func(db *sql.DB) {
		var datas []model.TaskRes
		rows, err := db.Query(`SELECT distinct(Taskres_pageurl), Taskres_pagetitle, Taskres_pagedate, Taskres_pagepath FROM taskres where Taskres_pagetitle like '%` + w + `%' or Taskres_pagetext like '%` + w + `%'`)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"iserror": true,
				"msg":     "搜索失败!请查看日志",
			})
			return
		}
		for rows.Next() {
			var tmp model.TaskRes
			err = rows.Scan(&tmp.Taskres_pageurl, &tmp.Taskres_pagetitle, &tmp.Taskres_pagedate, &tmp.Taskres_pagepath)
			//转一下
			tmp.Taskres_pagetext = html.UnescapeString(tmp.Taskres_pagetext)
			if err != nil {
				log.Println(err.Error())
			}
			datas = append(datas, tmp)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"iserror": true,
				"msg":     "搜索失败!",
			})
		} else {
			res, _ := json.Marshal(datas)
			c.JSON(http.StatusOK, gin.H{
				"iserror": false,
				"msg":     "搜索成功",
				"data":    string(res),
			})
		}
	})

}
