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

// CommentAction 发起评论
func CommentAction(ctx *gin.Context) {
	// 0. 获取信息
	var commentInfo constant.CommentActionParams
	err := ctx.ShouldBindQuery(&commentInfo)
	if err != nil {
		log.Errorf("commentInfo ShouldBindQuery failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	userID, _ := ctx.Get("UserID")
	commentInfo.UserID = userID.(int64)
	// todo 删除评论

	// 1. 更新comment表信息
	CommentActionRsp, err := utils.GetCommentSvrClient().CommentAction(ctx, &pb.CommentActionReq{
		UserId:      commentInfo.UserID,
		VideoId:     commentInfo.VideoId,
		CommentId:   commentInfo.CommentId,
		CommentText: commentInfo.CommentText,
		ActionType:  commentInfo.ActionType,
	})
	if err != nil {
		log.Errorf("CommentAction failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 2. 更新video表 comment_count
	_, err = utils.GetVideoSvrClient().UpdateCommentCount(ctx, &pb.UpdateCommentCountReq{
		VideoId:    commentInfo.VideoId,
		ActionType: commentInfo.ActionType,
	})
	if err != nil {
		log.Errorf("UpdateCommentCount failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 3. 返回响应
	// 3-1 评论者的详细信息
	if commentInfo.ActionType == 1 {
		commenterInfoRsp, err := utils.GetUserSvrClient().GetUserInfo(ctx, &pb.GetUserInfoRequest{
			UserId: commentInfo.UserID,
		})
		if err != nil {
			log.Errorf("GetUserInfo failed: %v", err)
			response.Fail(ctx, err.Error(), nil)
			return
		}
		CommentActionRsp.Comment.User = commenterInfoRsp.User
	}

	log.Infof("user%v comment on %v: %v success...", commentInfo.UserID, commentInfo.VideoId, commentInfo.CommentText)
	response.Success(ctx, "success", CommentActionRsp)

}

// CommentList 查看评论
func CommentList(ctx *gin.Context) {
	id := ctx.Query("video_id")
	videoId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Errorf("videoId invalid: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 1. 根据video_id获取comment表中信息
	commentListRsp, err := utils.GetCommentSvrClient().GetCommentList(ctx, &pb.GetCommentListReq{
		VideoId: videoId,
	})
	if err != nil {
		log.Errorf("GetCommentList failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 2. 填充视频评论用户信息
	// 2-1 拿到评论用户ID
	userIdList := make([]int64, 0)
	for _, comment := range commentListRsp.CommentList {
		userIdList = append(userIdList, comment.User.Id)
	}
	// 2-2 查询具体信息
	commentUserInfoRsp, err := utils.GetUserSvrClient().GetUserInfoDict(ctx, &pb.GetUserInfoDictRequest{
		UserIdList: userIdList,
	})
	if err != nil {
		log.Errorf("GetUserInfoDict failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 2-3 填充信息
	userMap := commentUserInfoRsp.UserInfoDict
	for _, comment := range commentListRsp.CommentList {
		comment.User = &pb.UserInfo{
			Id:              comment.User.Id,
			Name:            userMap[comment.User.Id].Name,
			Avatar:          userMap[comment.User.Id].Avatar,
			FollowCount:     userMap[comment.User.Id].FollowCount,
			FollowerCount:   userMap[comment.User.Id].FollowerCount,
			IsFollow:        userMap[comment.User.Id].IsFollow,
			BackgroundImage: userMap[comment.User.Id].BackgroundImage,
			Signature:       userMap[comment.User.Id].Signature,
			TotalFavorited:  userMap[comment.User.Id].TotalFavorited,
			FavoriteCount:   userMap[comment.User.Id].FavoriteCount,
		}
	}

	log.Infof("get CommentList success")
	response.Success(ctx, "success", commentListRsp)
}
