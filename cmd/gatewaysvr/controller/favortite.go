package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/constant"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

// FavoriteAction 点赞/取消赞
func FavoriteAction(ctx *gin.Context) {
	var favInfo constant.FavoriteActionParams
	err := ctx.ShouldBindQuery(&favInfo)
	if err != nil {
		log.Errorf("favInfo ShouldBindQuery failed: %v", err)
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

// FavoriteList 用户喜欢的视频
func FavoriteList(ctx *gin.Context) {
	userID, _ := ctx.Get("UserID")
	if userID == int64(-1) {
		log.Infof("login in first...")
		response.Fail(ctx, constant.ErrorNotLogin, nil)
		return
	}

	// 1. 根据userID,查favorite表中的videoIdlist
	favoriteVideoIdListRsp, err := utils.GetFavoriteSvrClient().GetFavoriteVideoIdList(ctx, &pb.GetFavoriteVideoIdListReq{
		UserId: userID.(int64),
	})
	if err != nil {
		log.Errorf("GetFavoriteVideoIdList failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 2. 根据videoId查询视频信息
	authorIdListRsp, err := utils.GetVideoSvrClient().GetVideoInfoList(ctx, &pb.GetVideoInfoListReq{
		VideoId: favoriteVideoIdListRsp.VideoIdList,
	})
	if err != nil {
		log.Errorf("GetVideoInfoList failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	var authorIdList []int64
	for _, v := range authorIdListRsp.VideoInfoList {
		authorIdList = append(authorIdList, v.AuthorId)
	}

	// 3.查询视频作者的信息
	videoAuthorInfoRep, err := utils.GetUserSvrClient().GetUserInfoDict(ctx, &pb.GetUserInfoDictRequest{
		UserIdList: authorIdList,
	})
	if err != nil {
		log.Errorf("GetVideoAuthorInfo failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	userMap := videoAuthorInfoRep.UserInfoDict
	// 4. 填充响应
	var rsp = &pb.FavoriteVideoListRsp{
		VideoList: make([]*pb.Video, 0),
	}

	for _, v := range authorIdListRsp.VideoInfoList {
		rsp.VideoList = append(rsp.VideoList, &pb.Video{
			Id:            v.Id,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
			Author: &pb.UserInfo{
				Id:              v.AuthorId,
				Name:            userMap[v.AuthorId].Name,
				Avatar:          userMap[v.AuthorId].Avatar,
				FollowCount:     userMap[v.AuthorId].FollowCount,
				FollowerCount:   userMap[v.AuthorId].FollowerCount,
				IsFollow:        userMap[v.AuthorId].IsFollow,
				BackgroundImage: userMap[v.AuthorId].BackgroundImage,
				Signature:       userMap[v.AuthorId].Signature,
				TotalFavorited:  userMap[v.AuthorId].TotalFavorited,
				FavoriteCount:   userMap[v.AuthorId].FavoriteCount,
			},
		})
	}

	log.Infof("get user: %v favoriteList", userID)
	response.Success(ctx, "success", rsp)
}
