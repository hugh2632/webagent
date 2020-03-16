package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"webagent/biz"
)

func CreateTask(c *gin.Context) {
	webid, ok := c.GetPostForm("webid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"iserror": true,
			"msg":     "网站id有误!",
		})
		return
	}
	onlyfirst, ok := c.GetPostForm("onlyfirst")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"iserror": true,
			"msg":     "缺少是否只扫主页参数!",
		})
		return
	}
	rebuild, ok := c.GetPostForm("rebuild")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"iserror": true,
			"msg":     "缺少是否重扫的参数!",
		})
		return
	}
	taskid, err := biz.TaskNewTask(webid, onlyfirst, rebuild)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"iserror": true,
			"msg":     err.Error(),
		})
		return
	} else {

		c.JSON(http.StatusOK, gin.H{
			"taskid": strconv.FormatUint(taskid, 10),
		})
	}
}

func RunTask(c *gin.Context) {
	var d, _ = c.GetPostForm("id")
	i64, _ := strconv.ParseUint(d, 10, 64)
	err := biz.TaskRun(i64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"iserror": true,
			"msg":     "发起扫描失败!" + err.Error(),
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"iserror": false,
			"msg":     "扫描成功",
		})
	}
}

func GetTaskRes(c *gin.Context) {
	var d, _ = c.GetPostForm("id")
	i64 , _ := strconv.ParseUint(d, 10, 64)
	info , res, err := biz.TaskGetRes(i64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"iserror": true,
			"msg":     "获取任务结果失败!",
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"iserror": false,
			"info":    info,
			"result":  res,
		})
	}
}

func TaskListSite(c *gin.Context) {
	res, err := biz.TaskListSite()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"iserror": true,
			"msg":     "获取网站列表失败!" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"iserror": false,
			"msg":     "扫描成功",
			"data" : res,
		})
	}
}

