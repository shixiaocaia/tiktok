package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/controller"
	jwtoken "github.com/shixiaocaia/tiktok/cmd/gatewaysvr/middleware"
)

func FavoriteRoutes(r *gin.RouterGroup) {
	favorite := r.Group("/favorite/", jwtoken.JWTAuthMiddleware())
	{
		favorite.POST("/action/", controller.FavoriteAction)
		favorite.GET("/list/", controller.FavoriteList)
	}
}
