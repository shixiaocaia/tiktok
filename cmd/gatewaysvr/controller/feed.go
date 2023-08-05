package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"strconv"
)

// Feed 视频流
func Feed(ctx *gin.Context) {
	// 从URL中读取时间戳转换为int64
	currentTime, err := strconv.ParseInt(ctx.Query("latest_time"), 10, 64)
	if err != nil || currentTime == int64(0) {
		currentTime = utils.GetCurrentTime()
	}

	// todo 根据token验证当前是登录状态还是未登录状态
	// 如果是登录要显示视频的参数，是否点赞，是否关注等状态
	userID, _ := ctx.Get("UserID")
	var UserID int64
	UserID = userID.(int64)
	// 未登录用户 UserID = 0

	// 调用rpc请求获取多个视频
	feedListResponse, err := utils.GetVideoSvrClient().GetFeedList(ctx, &pb.GetFeedListRequest{
		LatestTime: currentTime,
		UserId:     UserID,
	})
	// todo 视频顺序按照发布时间倒叙

	var authorIdList = make([]int64, 0)
	var followUintList = make([]*pb.FollowUint, 0)
	var favoriteUnitList = make([]*pb.FavoriteUnit, 0)

	for _, video := range feedListResponse.VideoList {
		// 视频作者ID
		authorIdList = append(authorIdList, video.AuthorId)
		// 视频关注
		followUintList = append(followUintList, &pb.FollowUint{
			SelfUserId: UserID,
			ToUserId:   video.AuthorId,
		})
		// 视频点赞
		favoriteUnitList = append(favoriteUnitList, &pb.FavoriteUnit{
			UserId:  UserID,
			VideoId: video.Id,
		})
	}
	// 根据视频的作者id，去查作者信息
	videoAuthorInfoRep, err := utils.GetUserSvrClient().GetUserInfoDict(ctx, &pb.GetUserInfoDictRequest{
		UserIdList: authorIdList,
	})
	if err != nil {
		log.Error("GetVideoAuthorInfo err...", err)
		response.Fail(ctx, fmt.Sprintf("GetUserSvrClient GetUserInfoDict err %v", err.Error()), nil)
		return
	}
	// log.Debugf("videoAuthorInfo: ", videoAuthorInfoRep.UserInfoDict)

	// todo 视频点赞

	// todo 关注用户

	// 填充响应返回
	var resp = &pb.DouyinFeedResponse{
		VideoList: make([]*pb.Video, 0),
		NextTime:  feedListResponse.NextTime,
	}
	for _, video := range feedListResponse.VideoList {
		videoRep := &pb.Video{
			Id:            video.Id,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		}
		// 作者详细信息
		videoRep.Author = videoAuthorInfoRep.UserInfoDict[video.AuthorId]
		// 登录用户，更新点赞和关注信息
		if UserID != 0 {

		}
		resp.VideoList = append(resp.VideoList, videoRep)
	}

	response.Success(ctx, "success", resp)
}
