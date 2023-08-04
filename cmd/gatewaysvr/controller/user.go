package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

func UserLogin(ctx *gin.Context) {
	userName := ctx.Query("username")
	passWord := ctx.Query("password")

	if len(userName) > 32 || len(passWord) > 32 {
		response.Fail(ctx, "username or password invalid", nil)
		return
	}

	resp, err := utils.GetUserSvrClient().CheckPassWord(ctx, &pb.CheckPassWordRequest{
		Username: userName,
		Password: passWord,
	})

	if err != nil {
		log.Error("Login failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	log.Info("login success")
	response.Success(ctx, "success", resp)
}

func UserRegister(ctx *gin.Context) {
	// 获取POST请求的userName, passWord
	userName := ctx.Query("username")
	passWord := ctx.Query("password")

	// 用户名和密码最长32个字符
	if len(userName) > 32 || len(passWord) > 32 {
		response.Fail(ctx, "username or password invalid", nil)
		return
	}

	// 调用userSvr的Register服务
	resp, err := utils.GetUserSvrClient().Register(ctx, &pb.RegisterRequest{
		Username: userName,
		Password: passWord,
	})

	if err != nil {
		log.Error("Register failed ", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	log.Info("register success")
	response.Success(ctx, "success", resp)
}

func GetUserInfo(ctx *gin.Context) {

}
