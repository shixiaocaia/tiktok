package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/controller"
	jwtoken "github.com/shixiaocaia/tiktok/cmd/gatewaysvr/middleware"
)

func CommentRoutes(r *gin.RouterGroup) {
	comment := r.Group("/comment/")
	{
		comment.POST("/action/", jwtoken.JWTAuthMiddleware(), controller.CommentAction)
		comment.GET("/list/", jwtoken.JWTWithoutAuthMiddleware(), controller.CommentList)
	}

}
