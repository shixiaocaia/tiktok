package service

import (
	"context"
	"github.com/shixiaocaia/tiktok/cmd/favoritesvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/favoritesvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"strconv"
)

type FavoriteService struct {
	pb.UnimplementedFavoriteServiceServer
}

func (u FavoriteService) FavoriteAction(ctx context.Context, req *pb.FavoriteActionReq) (*pb.FavoriteActionRsp, error) {
	// 判断操作 1/2
	var err error
	if req.ActionType == 1 {
		// 添加记录
		log.Infof("user: %v like video :%v", req.UserId, req.VideoId)
		err = dao.LikeAction(req.UserId, req.VideoId)
	} else {
		// 删除记录
		log.Infof("user: %v dislike video :%v", req.UserId, req.VideoId)
		err = dao.DislikeAction(req.UserId, req.VideoId)
	}
	if err != nil {
		log.Errorf("FavoriteAction failed: %v", err)
		return nil, err
	}

	rsp := &pb.FavoriteActionRsp{}
	return rsp, nil
}

func (u FavoriteService) IsFavoriteVideoDict(ctx context.Context, req *pb.IsFavoriteVideoDictReq) (*pb.IsFavoriteVideoDictRsp, error) {
	isFavoriteDict := make(map[string]bool)
	for _, unit := range req.FavoriteUnitList {
		isFavorite, err := dao.IsFavoriteVideo(unit.UserId, unit.VideoId)
		if err != nil {
			log.Errorf("IsFavoriteVideo failed: %v", err)
			return nil, err
		}
		isFavoriteKey := strconv.FormatInt(unit.UserId, 10) + "_" + strconv.FormatInt(unit.VideoId, 10)
		isFavoriteDict[isFavoriteKey] = isFavorite
	}
	return &pb.IsFavoriteVideoDictRsp{IsFavoriteDict: isFavoriteDict}, nil
}

func (u FavoriteService) GetFavoriteVideoIdList(ctx context.Context, req *pb.GetFavoriteVideoIdListReq) (*pb.GetFavoriteVideoIdListRsp, error) {
	userID := req.UserId
	videoIdList, err := dao.GetFavoriteVideoIdList(userID)
	if err != nil {
		log.Errorf("GetFavoriteVideoIdList failed: %v", err)
		return nil, err
	}
	rsp := &pb.GetFavoriteVideoIdListRsp{
		VideoIdList: videoIdList,
	}

	return rsp, nil
}
