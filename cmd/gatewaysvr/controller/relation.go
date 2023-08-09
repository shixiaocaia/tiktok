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

// RelationAction 关注操作
func RelationAction(ctx *gin.Context) {
	var relationInfo constant.RelationActionParams
	err := ctx.ShouldBindQuery(&relationInfo)
	if err != nil {
		log.Errorf("get relationInfo failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	userId, _ := ctx.Get("UserID")
	relationInfo.UserID = userId.(int64)

	// 1. 更新relation记录
	_, err = utils.GetRelationSvrClient().RelationAction(ctx, &pb.RelationActionReq{
		SelfUserId: relationInfo.UserID,
		ToUserId:   relationInfo.ToUserID,
		ActionType: relationInfo.ActionType,
	})
	if err != nil {
		log.Errorf("RelationAction failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 2. 更新User中follow，follower字段
	_, err = utils.GetUserSvrClient().UpdateUserFollowCount(ctx, &pb.UpdateUserFollowCountReq{
		UserId:     relationInfo.UserID,
		ActionType: relationInfo.ActionType,
	})
	if err != nil {
		log.Errorf("UpdateUserFollowCount failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	_, err = utils.GetUserSvrClient().UpdateUserFollowerCount(ctx, &pb.UpdateUserFollowerCountReq{
		UserId:     relationInfo.ToUserID,
		ActionType: relationInfo.ActionType,
	})
	if err != nil {
		log.Errorf("UpdateUserFollowerCount failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 3. 返回响应
	log.Infof("user%v follow/not follow user%v", relationInfo.UserID, relationInfo.ToUserID)
	rsp := &pb.RelationActionRsp{}
	response.Success(ctx, "success", rsp)
}

// FollowList 获取关注列表
func FollowList(ctx *gin.Context) {
	userId, _ := ctx.Get("UserID")

	// 1. 获取自己的ID
	followListRsp, err := utils.RelationClient.GetRelationFollowList(ctx, &pb.GetRelationFollowListReq{
		UserId: userId.(int64),
	})
	if err != nil {
		log.Errorf("GetRelationFollowList failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 2. 根据id查询关注信息
	followIdList := make([]int64, 0)
	for _, follow := range followListRsp.UserList {
		followIdList = append(followIdList, follow.Id)
	}

	followInfoRsp, err := utils.GetUserSvrClient().GetUserInfoDict(ctx, &pb.GetUserInfoDictRequest{
		UserIdList: followIdList,
	})
	if err != nil {
		log.Errorf("GetUserInfoDict failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 3. 填充具体信息
	InfoMap := followInfoRsp.UserInfoDict
	followUser := make([]*pb.UserInfo, 0)
	for _, follow := range followListRsp.UserList {
		Info := InfoMap[follow.Id]
		followUser = append(followUser, &pb.UserInfo{
			Id:              Info.Id,
			Name:            Info.Name,
			FollowCount:     Info.FollowCount,
			FollowerCount:   Info.FollowerCount,
			IsFollow:        true,
			Avatar:          Info.Avatar,
			BackgroundImage: Info.BackgroundImage,
			Signature:       Info.Signature,
			TotalFavorited:  Info.TotalFavorited,
			WorkCount:       Info.WorkCount,
			FavoriteCount:   Info.FavoriteCount,
		})
	}

	log.Infof("get user %v follow", userId)
	response.Success(ctx, "success", &pb.GetRelationFollowListRsp{
		UserList: followUser,
	})
}

// FollowerList 获取粉丝列表
func FollowerList(ctx *gin.Context) {
	Id := ctx.Query("user_id")
	userId, err := strconv.ParseInt(Id, 10, 64)
	// 1. 获取粉丝的id
	followerListRsp, err := utils.RelationClient.GetRelationFollowerList(ctx, &pb.GetRelationFollowerListReq{
		UserId: userId,
	})
	log.Debugf("his follower: %v", followerListRsp.UserList)
	if err != nil {
		log.Errorf("GetRelationFollowerList failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 2. 根据id查询粉丝信息
	followerIdList := make([]int64, 0)
	for _, follower := range followerListRsp.UserList {
		followerIdList = append(followerIdList, follower.Id)
	}
	followInfoRsp, err := utils.GetUserSvrClient().GetUserInfoDict(ctx, &pb.GetUserInfoDictRequest{
		UserIdList: followerIdList,
	})
	if err != nil {
		log.Errorf("GetUserInfoDict failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 3. 获取自己关注哪些人
	SelfId, _ := ctx.Get("UserID")
	FollowDic, err := utils.GetRelationSvrClient().IsFollowDict(ctx, &pb.IsFollowDictReq{
		UserId: SelfId.(int64),
	})
	if err != nil {
		log.Errorf("IsFollowDict failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 4. 填充具体信息
	InfoMap := followInfoRsp.UserInfoDict
	followMap := FollowDic.IsFollowDict
	log.Debugf("InfoMap: %v", followMap)
	followerUser := make([]*pb.UserInfo, 0)
	for _, follower := range followerListRsp.UserList {
		Info := InfoMap[follower.Id]
		log.Debugf("is: %v", followMap[follower.Id])
		followerUser = append(followerUser, &pb.UserInfo{
			Id:            Info.Id,
			Name:          Info.Name,
			FollowCount:   Info.FollowCount,
			FollowerCount: Info.FollowerCount,
			// todo是否关注了这名粉丝
			IsFollow:        followMap[follower.Id],
			Avatar:          Info.Avatar,
			BackgroundImage: Info.BackgroundImage,
			Signature:       Info.Signature,
			TotalFavorited:  Info.TotalFavorited,
			WorkCount:       Info.WorkCount,
			FavoriteCount:   Info.FavoriteCount,
		})
	}

	log.Infof("get user %v follower", userId)
	response.Success(ctx, "success", &pb.GetRelationFollowerListRsp{
		UserList: followerUser,
	})
}

func FriendList(ctx *gin.Context) {
	Id := ctx.Query("user_id")
	userId, err := strconv.ParseInt(Id, 10, 64)
	// 1. 获取粉丝的id
	followerListRsp, err := utils.RelationClient.GetRelationFollowerList(ctx, &pb.GetRelationFollowerListReq{
		UserId: userId,
	})

	if err != nil {
		log.Errorf("GetRelationFollowerList failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 2. 根据id查询粉丝信息
	followerIdList := make([]int64, 0)
	for _, follower := range followerListRsp.UserList {
		followerIdList = append(followerIdList, follower.Id)
	}
	followInfoRsp, err := utils.GetUserSvrClient().GetUserInfoDict(ctx, &pb.GetUserInfoDictRequest{
		UserIdList: followerIdList,
	})
	if err != nil {
		log.Errorf("GetUserInfoDict failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	// 3. 获取自己关注哪些人
	SelfId, _ := ctx.Get("UserID")
	FollowDic, err := utils.GetRelationSvrClient().IsFollowDict(ctx, &pb.IsFollowDictReq{
		UserId: SelfId.(int64),
	})
	if err != nil {
		log.Errorf("IsFollowDict failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 4. 填充具体信息
	InfoMap := followInfoRsp.UserInfoDict
	followMap := FollowDic.IsFollowDict
	followerUser := make([]*pb.UserInfo, 0)
	for _, follower := range followerListRsp.UserList {
		Info := InfoMap[follower.Id]
		followerUser = append(followerUser, &pb.UserInfo{
			Id:            Info.Id,
			Name:          Info.Name,
			FollowCount:   Info.FollowCount,
			FollowerCount: Info.FollowerCount,
			// todo是否关注了这名粉丝
			IsFollow:        followMap[follower.Id],
			Avatar:          Info.Avatar,
			BackgroundImage: Info.BackgroundImage,
			Signature:       Info.Signature,
			TotalFavorited:  Info.TotalFavorited,
			WorkCount:       Info.WorkCount,
			FavoriteCount:   Info.FavoriteCount,
		})
	}

	log.Infof("get user %v follower", userId)
	response.Success(ctx, "success", &pb.GetRelationFollowerListRsp{
		UserList: followerUser,
	})
}
