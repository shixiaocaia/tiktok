package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/controller"
)

func UserRoutes(r *gin.RouterGroup) {
	user := r.Group("user")
	{
		user.POST("/login/", controller.UserLogin)
		user.GET("/", controller.GetUserInfo)

		user.POST("/register/", controller.UserRegister)
	}

}
