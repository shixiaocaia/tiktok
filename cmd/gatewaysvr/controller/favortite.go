package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

type FavActionParams struct {
	UserID     int64
	VideoId    int64 `form:"video_id" binding:"required"`
	ActionType int64 `form:"action_type" binding:"required,oneof=1 2"`
}

func FavoriteAction(ctx *gin.Context) {
	var favInfo FavActionParams
	err := ctx.ShouldBindQuery(&favInfo)
	if err != nil {
		log.Errorf("FavoriteAction ShouldBindQuery failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	id, _ := ctx.Get("UserID")
	favInfo.UserID = id.(int64)

	// 更新favorite表
	_, err = utils.GetFavoriteSvrClient().FavoriteAction(ctx, &pb.FavoriteActionReq{
		UserId:     favInfo.UserID,
		VideoId:    favInfo.VideoId,
		ActionType: favInfo.ActionType,
	})
	if err != nil {
		log.Errorf("GetFavoriteSvrClient().FavoriteAction failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 1.更新 video表中 favorite_count
	_, err = utils.GetVideoSvrClient().UpdateFavoriteCount(ctx, &pb.UpdateFavoriteCountReq{
		VideoId:    favInfo.VideoId,
		ActionType: favInfo.ActionType,
	})
	if err != nil {
		log.Errorf("UpdateFavoriteCount failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 2.更新 user表中 favorite_count
	// 2-1. 先查根据video_id 查author_id
	// todo 能不能简化
	videoInfoRsp, err := utils.GetVideoSvrClient().GetVideoInfoList(ctx, &pb.GetVideoInfoListReq{
		VideoId: []int64{favInfo.VideoId},
	})
	if err != nil {
		log.Errorf("GetVideoInfoList failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 2-2. 根据author_id 更新user表 获赞数total_favorited
	_, err = utils.GetUserSvrClient().UpdateUserFavoritedCount(ctx, &pb.UpdateUserFavoritedCountReq{
		UserId:     videoInfoRsp.VideoInfoList[0].AuthorId,
		ActionType: favInfo.ActionType,
	})
	if err != nil {
		log.Errorf("UpdateUserFavoritedCount failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 3. 更新user表favorite_count
	_, err = utils.GetUserSvrClient().UpdateUserFavoriteCount(ctx, &pb.UpdateUserFavoriteCountReq{
		UserId:     favInfo.UserID,
		ActionType: favInfo.ActionType,
	})
	if err != nil {
		log.Errorf("UpdateUserFavoriteCount failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	log.Infof("user: %v like/dislike video: %v", favInfo.UserID, favInfo.VideoId)
	response.Success(ctx, "success", nil)
}

func FavoriteList(ctx *gin.Context) {

}
