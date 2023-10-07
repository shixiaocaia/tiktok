package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/controller"
	jwtoken "github.com/shixiaocaia/tiktok/cmd/gatewaysvr/middleware"
)

func RalationRoutes(c *gin.RouterGroup) {
	relation := c.Group("relation")
	{
		relation.POST("/action/", jwtoken.JWTAuthMiddleware(), controller.RelationAction)
		relation.GET("/follow/list/", jwtoken.JWTAuthMiddleware(), controller.FollowList)
		relation.GET("/follower/list/", jwtoken.JWTAuthMiddleware(), controller.FollowerList)
		relation.GET("/friend/list/", jwtoken.JWTAuthMiddleware(), controller.FriendList)
	}
}
