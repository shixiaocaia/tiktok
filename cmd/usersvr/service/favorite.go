package service

import (
	"context"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

// UpdateUserFavoritedCount 更新获赞数
func (u UserService) UpdateUserFavoritedCount(ctx context.Context, req *pb.UpdateUserFavoritedCountReq) (*pb.UpdateUserFavoritedCountRsp, error) {
	err := dao.UpdateFavorited(req.UserId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateFavorited failed: %v", err)
		return nil, err
	}

	return &pb.UpdateUserFavoritedCountRsp{}, nil
}

// UpdateUserFavoriteCount 更新喜欢数
func (u UserService) UpdateUserFavoriteCount(ctx context.Context, req *pb.UpdateUserFavoriteCountReq) (*pb.UpdateUserFavoriteCountRsp, error) {
	err := dao.UpdateFavorite(req.UserId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateFavorite failed: %v", err)
		return nil, err
	}

	return &pb.UpdateUserFavoriteCountRsp{}, nil
}
