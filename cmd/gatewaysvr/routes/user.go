package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/controller"
	jwtoken "github.com/shixiaocaia/tiktok/cmd/gatewaysvr/middleware"
)

func UserRoutes(r *gin.RouterGroup) {
	user := r.Group("user")
	{
		user.POST("/register/", controller.UserRegister)
		user.POST("/login/", controller.UserLogin)
		user.GET("/", jwtoken.JWTAuthMiddleware(), controller.GetUserInfo)
	}

}
