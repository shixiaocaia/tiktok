package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"videosvr/config"
	"videosvr/log"
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
