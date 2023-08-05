package service

import (
	"context"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"time"
)

type VideoService struct {
	pb.UnimplementedVideoServiceServer
}

// GetFeedList 获取一组视频基本信息
func (u VideoService) GetFeedList(ctx context.Context, req *pb.GetFeedListRequest) (*pb.GetFeedListResponse, error) {
	// 返回以时间戳为止的一组视频
	videoList, err := dao.GetVideoListByFeed(req.LatestTime)
	if err != nil {
		log.Error("GetVideoListByFeed failed")
		return nil, err
	}
	// 返回下一批视频的最新时间
	nextTime := time.Now().UnixNano() / 1e6
	if len(videoList) == 20 {
		nextTime = videoList[len(videoList)-1].PublishTime
	}
	//
	var VideoInfoList []*pb.VideoInfo
	for _, video := range videoList {
		VideoInfoList = append(VideoInfoList, &pb.VideoInfo{
			Id:            video.Id,
			AuthorId:      video.AuthorId,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    false, // 是否喜欢，在gateway处理
			Title:         video.Title,
		})
	}

	resp := &pb.GetFeedListResponse{
		VideoList: VideoInfoList,
		NextTime:  nextTime,
	}
	return resp, nil
}
