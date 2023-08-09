package service

import (
	"context"
	"github.com/shixiaocaia/tiktok/cmd/relationsvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/relationsvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

type RelationService struct {
	pb.UnimplementedRelationServiceServer
}

func (u RelationService) RelationAction(ctx context.Context, req *pb.RelationActionReq) (*pb.RelationActionRsp, error) {
	var err error
	if req.ActionType == 1 {
		err = dao.FollowAdd(req.ToUserId, req.SelfUserId)
		if err != nil {
			log.Errorf("FollowAdd failed: %v", err)
			return nil, err
		}
	} else {
		err = dao.FollowDel(req.ToUserId, req.SelfUserId)
		if err != nil {
			log.Errorf("FollowDel failed: %v", err)
			return nil, err
		}
	}
	return &pb.RelationActionRsp{}, nil
}

func (u RelationService) GetRelationFollowList(ctx context.Context, req *pb.GetRelationFollowListReq) (*pb.GetRelationFollowListRsp, error) {
	userId := req.UserId
	followIdList, err := dao.GetFollowIdList(userId)
	if err != nil {
		log.Errorf("GetFollowIdList failed: %v", err)
		return nil, err
	}

	var rsp = &pb.GetRelationFollowListRsp{
		UserList: make([]*pb.UserInfo, 0),
	}
	for _, follow := range followIdList {
		rsp.UserList = append(rsp.UserList, &pb.UserInfo{
			Id: follow.Follow,
		})
	}
	return rsp, nil
}

func (u RelationService) GetRelationFollowerList(ctx context.Context, req *pb.GetRelationFollowerListReq) (*pb.GetRelationFollowerListRsp, error) {
	userId := req.UserId
	followIdList, err := dao.GetFollowerIdList(userId)
	if err != nil {
		log.Errorf("GetFollowIdList failed: %v", err)
		return nil, err
	}

	var rsp = &pb.GetRelationFollowerListRsp{
		UserList: make([]*pb.UserInfo, 0),
	}
	for _, follow := range followIdList {
		rsp.UserList = append(rsp.UserList, &pb.UserInfo{
			Id: follow.Follower,
		})
	}
	return rsp, nil

}

func (u RelationService) IsFollowDict(ctx context.Context, req *pb.IsFollowDictReq) (*pb.IsFollowDictRsp, error) {
	userId := req.UserId
	followIdList, err := dao.GetFollowIdList(userId)
	if err != nil {
		log.Errorf("GetFollowIdList failed: %v", err)
		return nil, err
	}

	rsp := &pb.IsFollowDictRsp{IsFollowDict: make(map[int64]bool)}
	for _, follow := range followIdList {
		rsp.IsFollowDict[follow.Follow] = true
	}
	return rsp, nil
}

func (u RelationService) IsFriendList(ctx context.Context, req *pb.IsFriendListReq) (*pb.IsFriendListRsp, error) {
	friendList := make([]int64, 0)
	for _, follower := range req.UserList {
		flag, err := dao.IsFriend(follower, req.UserId)
		if err != nil {
			log.Errorf("dao.IsFriend failed: %v", err)
			return nil, err
		}
		if flag == true {
			friendList = append(friendList, follower)
		}
	}
	return &pb.IsFriendListRsp{
		UserList: friendList,
	}, nil
}
