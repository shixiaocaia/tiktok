package service

import (
	"context"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

func (u UserService) UpdateUserFollowCount(ctx context.Context, req *pb.UpdateUserFollowCountReq) (*pb.UpdateUserFollowCountRsp, error) {
	log.Debugf("req %v", req)
	err := dao.UpdateFollowCount(req.UserId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateFollowCount failed: %v", err)
		return nil, err
	}
	return &pb.UpdateUserFollowCountRsp{}, nil
}

func (u UserService) UpdateUserFollowerCount(ctx context.Context, req *pb.UpdateUserFollowerCountReq) (*pb.UpdateUserFollowerCountRsp, error) {
	err := dao.UpdateFollowerCount(req.UserId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateFollowerCount failed: %v", err)
		return nil, err
	}
	return &pb.UpdateUserFollowerCountRsp{}, nil
}
