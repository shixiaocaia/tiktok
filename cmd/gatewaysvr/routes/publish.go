package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/controller"
	jwtoken "github.com/shixiaocaia/tiktok/cmd/gatewaysvr/middleware"
)

func PublishVideoRoutes(r *gin.RouterGroup) {
	// JWT
	publish := r.Group("publish")
	{
		publish.POST("/action/", jwtoken.JWTAuthMiddleware(), controller.PublishAction)
		publish.GET("/list/", jwtoken.JWTWithoutAuthMiddleware(), controller.GetPublishList)
	}

}
