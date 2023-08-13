package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/config"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/constant"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"path/filepath"
)

// PublishAction 发布视频
func PublishAction(ctx *gin.Context) {
	userID, _ := ctx.Get("UserID")
	title := ctx.PostForm("title")
	data, err := ctx.FormFile("data")
	if err != nil {
		log.Errorf("upload video failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	filename := filepath.Base(data.Filename)
	// 本地存放
	finalName := fmt.Sprintf("%s_%s", utils.RandomString(), filename)
	videoPath := config.GetGlobalConfig().VideoPath
	saveFile := filepath.Join(videoPath, finalName)
	if err := ctx.SaveUploadedFile(data, saveFile); err != nil {
		log.Errorf("UploadVideo failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 上传视频, 截取视频封面并上传minio
	_, err = utils.GetVideoSvrClient().PublishVideo(ctx, &pb.PublishVideoRequest{
		UserId:   userID.(int64),
		Title:    title,
		SaveFile: saveFile,
	})
	if err != nil {
		log.Errorf("PublishVideo failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 更新用户作品数
	_, err = utils.GetUserSvrClient().UpdateWorkCount(ctx, &pb.UpdateUserWorkCountReq{
		UserId: userID.(int64),
	})
	if err != nil {
		log.Errorf("UpdateWorkCount failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	log.Info("publish video success")
	response.Success(ctx, "success", &pb.PublishVideoResponse{})
}

// GetPublishList 发布列表
func GetPublishList(ctx *gin.Context) {
	UserID, _ := ctx.Get("UserID")
	if UserID == int64(-1) {
		log.Infof("login in first...")
		response.Fail(ctx, constant.ErrorNotLogin, nil)
		return
	}
	// 获取视频列表
	getPublishList, err := utils.GetVideoSvrClient().GetPublishVideoList(ctx, &pb.GetPublishVideoListRequest{
		UserId: UserID.(int64),
	})
	if err != nil {
		log.Errorf("GetPublishVideoList failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 获取作者信息
	getUserInfo, err := utils.GetUserSvrClient().GetUserInfo(ctx, &pb.GetUserInfoRequest{
		UserId: UserID.(int64),
	})

	if err != nil {
		log.Errorf("GetUserInfo err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 合并视频列表和作者信息
	for _, video := range getPublishList.VideoList {
		video.Author = getUserInfo.User
	}

	log.Infof("get author: %v videos", UserID)
	response.Success(ctx, "success", getPublishList)
}
