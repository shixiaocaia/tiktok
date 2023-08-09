package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/controller"
	jwtoken "github.com/shixiaocaia/tiktok/cmd/gatewaysvr/middleware"
)

func MessageRoutes(c *gin.RouterGroup) {
	message := c.Group("message")
	{
		message.GET("/chat/", jwtoken.JWTAuthMiddleware(), controller.MessageChat)
		message.POST("/action/", jwtoken.JWTAuthMiddleware(), controller.MessageAction)
	}
}
