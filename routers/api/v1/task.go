package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"webagent/biz"
)

func TaskCreateTask(c *gin.Context) {
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
			"taskid": taskid,
		})
	}
}

func TaskRunTask(c *gin.Context) {
	var d, _ = c.GetPostForm("id")
	err := biz.TaskRun(d)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"iserror": true,
			"msg":     "发起扫描失败!" + err.Error(),
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"iserror": false,
			"msg":     "启动扫描成功",
		})
	}
}

func TaskGetTaskRes(c *gin.Context) {
	var d, _ = c.GetPostForm("id")
	info, res, err := biz.TaskGetRes(d)
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
			"msg":     "获取成功",
			"data":    res,
		})
	}
}

func TaskListTask(c *gin.Context) {
	order, _ := c.GetQuery("order")
	sort, _ := c.GetQuery("sort")
	pageindex, b := c.GetQuery("pageindex")
	rowcount, _ := c.GetQuery("rowcount")
	var indexNum = -1
	var countNum = -1
	if b {
		n, err := strconv.Atoi(pageindex)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"iserror": true,
				"msg":     "pageindex非数字",
			})
			return
		}
		indexNum = n
		cn, err := strconv.Atoi(rowcount)
		if err == nil {
			countNum = cn
		}
	}
	res , total , err := biz.TaskListTaskInfo(sort, order, indexNum, countNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"iserror": true,
			"msg":     "获取网站列表失败!" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"iserror": false,
			"msg":     "获取成功",
			"data":    res,
			"total":  total,
		})
	}
}