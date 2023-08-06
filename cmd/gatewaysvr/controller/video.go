package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/config"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"path/filepath"
)

// PublishAction 发布视频
func PublishAction(ctx *gin.Context) {
	// JWT鉴权后获得userid
	utils.GetRequestInfo(ctx)
	userID, _ := ctx.Get("UserID")
	title := ctx.PostForm("title")
	data, err := ctx.FormFile("data")

	log.Debugf("userID: %v, title: %v", userID, title)

	if err != nil {
		log.Errorf("upload video failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	filename := filepath.Base(data.Filename)

	// 文件名 + 视频存放路径 + 最终保存路径
	finalName := fmt.Sprintf("%s_%s", utils.RandomString(), filename)
	videoPath := config.GetGlobalConfig().VideoPath
	saveFile := filepath.Join(videoPath, finalName)

	log.Debugf("videoPath:%v", videoPath)

	if err := ctx.SaveUploadedFile(data, saveFile); err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

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
	response.Success(ctx, "success", &pb.PublishVideoResponse{})
}

// GetPublishList 发布列表
func GetPublishList(ctx *gin.Context) {
	tokenUserId, _ := ctx.Get("UserID")
	//id := ctx.Query("user_id")
	log.Debugf("tokenUserId: %v", tokenUserId)
	// 获取视频
	getPublishList, err := utils.GetVideoSvrClient().GetPublishVideoList(ctx, &pb.GetPublishVideoListRequest{
		UserId: tokenUserId.(int64),
	})
	log.Debugf("getPublishList: %v", getPublishList)
	if err != nil {
		log.Errorf("GetPublishVideoList failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 获取作者信息
	getUserInfo, err := utils.GetUserSvrClient().GetUserInfo(ctx, &pb.GetUserInfoRequest{
		UserId: tokenUserId.(int64),
	})

	if err != nil {
		log.Errorf("GetUserInfo err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	for _, video := range getPublishList.VideoList {
		video.Author = getUserInfo.User
	}

	log.Debugf("getPublishList: %v", getPublishList)
	response.Success(ctx, "success", getPublishList)

}
