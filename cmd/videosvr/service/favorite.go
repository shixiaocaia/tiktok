package service

import (
	"context"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

func (u VideoService) UpdateFavoriteCount(ctx context.Context, req *pb.UpdateFavoriteCountReq) (*pb.UpdateFavoriteCountRsp, error) {
	err := dao.UpdateFavorite(req.ActionType, req.VideoId)
	if err != nil {
		log.Errorf("UpdateFavoriteCount failed: %v", err)
		return nil, err
	}
	return &pb.UpdateFavoriteCountRsp{}, nil
}
