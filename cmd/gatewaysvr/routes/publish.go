package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/controller"
	jwtoken "github.com/shixiaocaia/tiktok/cmd/gatewaysvr/middleware"
)

func PublishVideoRoutes(r *gin.RouterGroup) {
	// JWT
	publish := r.Group("publish", jwtoken.JWTAuthMiddleware())
	{
		publish.POST("/action/", controller.PublishAction)
		publish.GET("/list/", controller.GetPublishList)
	}

}
