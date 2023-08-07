package service

import (
	"context"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/dao"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

func (u UserService) UpdateWorkCount(ctx context.Context, req *pb.UpdateUserWorkCountReq) (*pb.UpdateUserWorkCountRsp, error) {
	authorID := req.UserId
	err := dao.UpdateWorkCount(authorID)
	if err != nil {
		log.Errorf("UpdateWorkCount failed")
		return nil, err
	}
	return &pb.UpdateUserWorkCountRsp{}, nil
}
