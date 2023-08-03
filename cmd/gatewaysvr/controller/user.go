package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"go.uber.org/zap"
)

func UserLogin(ctx *gin.Context) {

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
		zap.L().Error("login error", zap.Error(err))
		response.Fail(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, "success", resp)
}

func GetUserInfo(ctx *gin.Context) {

}
