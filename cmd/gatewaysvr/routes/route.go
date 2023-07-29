package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/config"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"net/http"
)

func SetRoute() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": config.GetGlobalConfig().Ping,
		})
		log.Info(config.GetGlobalConfig().Ping)
	})

	return r
}
