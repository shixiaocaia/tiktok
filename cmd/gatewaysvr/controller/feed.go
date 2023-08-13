package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"strconv"
)

// Feed 视频流
func Feed(ctx *gin.Context) {
	// 1. 获取时间戳，判断有无token，从token提取用户ID
	currentTime, err := strconv.ParseInt(ctx.Query("latest_time"), 10, 64)
	if err != nil || currentTime == int64(0) {
		currentTime = utils.GetCurrentTime()
	}
	// 1-1. 未登录用户，从ctx读到-1，区分登录还是未登录
	userID, _ := ctx.Get("UserID")
	UserID := userID.(int64)

	// 2. 获取一组视频以及相关信息
	// 2-1根据时间戳拿到一组视频
	feedListResponse, err := utils.GetVideoSvrClient().GetFeedList(ctx, &pb.GetFeedListRequest{
		LatestTime: currentTime,
		UserId:     UserID,
	})
	var authorIdList = make([]int64, 0)
	var favoriteUnitList = make([]*pb.FavoriteUnit, 0)
	for _, video := range feedListResponse.VideoList {
		// 视频作者ID
		authorIdList = append(authorIdList, video.AuthorId)
		// 视频点赞信息
		favoriteUnitList = append(favoriteUnitList, &pb.FavoriteUnit{
			UserId:  UserID,
			VideoId: video.Id,
		})
	}

	// 2-2根据视频的作者id，去查作者信息
	videoAuthorInfoRep, err := utils.GetUserSvrClient().GetUserInfoDict(ctx, &pb.GetUserInfoDictRequest{
		UserIdList: authorIdList,
	})
	if err != nil {
		log.Errorf("GetVideoAuthorInfo failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	var videoFavoriteListRep *pb.IsFavoriteVideoDictRsp
	var FollowDic *pb.IsFollowDictRsp
	if UserID != -1 {
		// 2-3 登录用户，判断视频是否点赞
		videoFavoriteListRep, err = utils.GetFavoriteSvrClient().IsFavoriteVideoDict(ctx, &pb.IsFavoriteVideoDictReq{
			FavoriteUnitList: favoriteUnitList,
		})
		if err != nil {
			log.Errorf("IsFavoriteVideoDict failed: %v", err)
			response.Fail(ctx, err.Error(), nil)
			return
		}

		// 2-4 登录用户，判断视频作者是否关注
		SelfId, _ := ctx.Get("UserID")
		FollowDic, err = utils.GetRelationSvrClient().IsFollowDict(ctx, &pb.IsFollowDictReq{
			UserId: SelfId.(int64),
		})
		if err != nil {
			log.Errorf("IsFollowDict failed: %v", err)
			response.Fail(ctx, err.Error(), nil)
			return
		}
	}
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
			Title:         video.Title,
		}
		// 作者详细信息
		videoRep.Author = videoAuthorInfoRep.UserInfoDict[video.AuthorId]
		// 登录用户，更新点赞和关注信息
		if UserID != -1 {
			var favoriteUint = strconv.FormatInt(UserID, 10) + "_" + strconv.FormatInt(videoRep.Id, 10)
			videoRep.IsFavorite = videoFavoriteListRep.IsFavoriteDict[favoriteUint]
			videoRep.Author.IsFollow = FollowDic.IsFollowDict[video.AuthorId]
			// 忽略个人关注信息
			if video.AuthorId == UserID {
				videoRep.Author.IsFollow = true
			}
		}
		resp.VideoList = append(resp.VideoList, videoRep)
	}

	log.Infof("getFeed success")
	response.Success(ctx, "success", resp)
}
