package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/controller"
	jwtoken "github.com/shixiaocaia/tiktok/cmd/gatewaysvr/middleware"
)

func FavoriteRoutes(r *gin.RouterGroup) {
	favorite := r.Group("/favorite/")
	{
		favorite.POST("/action/", jwtoken.JWTAuthMiddleware(), controller.FavoriteAction)
		favorite.GET("/list/", jwtoken.JWTWithoutAuthMiddleware(), controller.FavoriteList)
	}
}
