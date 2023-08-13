package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/constant"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"strconv"
)

func UserLogin(ctx *gin.Context) {
	userName := ctx.Query("username")
	passWord := ctx.Query("password")

	// 1. 简单参数校验
	if len(userName) > 32 || len(passWord) > 32 {
		response.Fail(ctx, constant.InvalidUserInfo, nil)
		return
	}
	// 2. 调用userSvr的CheckPassWord服务
	resp, err := utils.GetUserSvrClient().CheckPassWord(ctx, &pb.CheckPassWordRequest{
		Username: userName,
		Password: passWord,
	})

	if err != nil {
		log.Error("Login failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	log.Info("login success...")
	response.Success(ctx, "success", resp)
}

func UserRegister(ctx *gin.Context) {
	// 获取POST请求的userName, passWord
	userName := ctx.Query("username")
	passWord := ctx.Query("password")

	// 用户名和密码最长32个字符
	if len(userName) > 32 || len(passWord) > 32 {
		response.Fail(ctx, constant.InvalidUserInfo, nil)
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
	log.Info("UserRegister success")
	response.Success(ctx, "success", resp)
}

func GetUserInfo(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	uids, ok := ctx.Get("UserID")
	if !ok {
		log.Error("cannot get uids from ctx")
		response.Fail(ctx, constant.ErrorToken, nil)
		return
	}
	uid := uids.(int64)

	// 1. 简单参数校验
	if strconv.FormatInt(uid, 10) != userId {
		log.Error("invalid uid")
		response.Fail(ctx, constant.ErrorToken, nil)
		return
	}

	// 2. 根据userID获取用户信息
	resp, err := utils.GetUserSvrClient().GetUserInfo(ctx, &pb.GetUserInfoRequest{
		UserId: uid,
	})
	if err != nil {
		log.Error("GetUserInfo failed ", err)
		response.Fail(ctx, constant.ErrorToken, nil)
		return
	}

	log.Info("getUserinfo success...")
	response.Success(ctx, "success", resp)
}
