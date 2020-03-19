package routers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"webagent/pkg/setting"
	"webagent/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(timeoutMiddleware(setting.ReadTimeout))

	gin.SetMode(setting.RunMode)
	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/search", v1.Search)
		apiv1.GET("/listsite", v1.TaskListSite)
		apiv1.POST("/runtask", v1.RunTask)
		apiv1.POST("/createtask", v1.CreateTask)
		apiv1.POST("/gettaskres", v1.GetTaskRes)
	}
	return r
}

func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer func() {
			if ctx.Err() == context.DeadlineExceeded {
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}
			cancel()
		}()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
