package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/config"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/controller"
	jwtoken "github.com/shixiaocaia/tiktok/cmd/gatewaysvr/middleware"
)

func SetRoute() *gin.Engine {
	if config.GetGlobalConfig().SvrConfig.Mode == gin.ReleaseMode {
		// gin设置成发布模式：gin不在终端输出日志
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()

	// test
	//r.GET("/ping", controller.Ping)
	//r.GET("/greet", controller.Greet)

	douyin := r.Group("/douyin/")
	{
		UserRoutes(douyin)
		PublishVideoRoutes(douyin)
		FavoriteRoutes(douyin)
		douyin.GET("/feed/", jwtoken.JWTWithoutAuthMiddleware(), controller.Feed)
	}
	return r
}
