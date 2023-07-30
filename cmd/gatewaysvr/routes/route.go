package routes

import (
	"gatewaysvr/config"
	"gatewaysvr/controller"
	"gatewaysvr/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetRoute() *gin.Engine {
	if config.GetGlobalConfig().Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": config.GetGlobalConfig().Ping,
		})
		log.Info(config.GetGlobalConfig().Ping)
	})

	douyin := r.Group("/douyin/")
	{

		douyin.GET("/feed/", controller.Feed)
	}

	return r
}
