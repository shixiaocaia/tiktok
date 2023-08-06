package service

import (
	"context"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/log"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/middleware/minio"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"strconv"
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

// PublishVideo 上传视频
func (u VideoService) PublishVideo(ctx context.Context, req *pb.PublishVideoRequest) (*pb.PublishVideoResponse, error) {
	// 连接到minio
	client := minio.GetMinio()
	playUrl, err := client.UploadFile("video", req.SaveFile, strconv.FormatInt(req.UserId, 10))
	if err != nil {
		log.Errorf("Minio UploadFile err: %v", err)
		return nil, err
	}

	log.Infof("save file: %v", req.SaveFile)

	// 生成视频封面
	imageFile, err := utils.GetImageFile(req.SaveFile)
	if err != nil {
		log.Errorf("ffmpeg getVideoPic failed: %v", err)
		return nil, err
	}

	coverUrl, err := client.UploadFile("pic", imageFile, strconv.FormatInt(req.UserId, 10))
	if err != nil {
		log.Errorf("minio upLoadPic failed: %v", err)
		return nil, err
	}

	log.Debugf("title: %v", req.Title)
	err = dao.InsertVideo(req.UserId, playUrl, coverUrl, req.Title)
	if err != nil {
		log.Errorf("InsertVideo failed: %v", err)
		return nil, err
	}
	return &pb.PublishVideoResponse{}, nil
}

func (u VideoService) GetPublishVideoList(ctx context.Context, req *pb.GetPublishVideoListRequest) (*pb.GetPublishVideoListResponse, error) {
	videos, err := dao.GetVideoListByAuthorID(req.UserId)
	if err != nil {
		log.Errorf("GetVideoListByAuthorID failed: %v", err)
		return nil, err
	}

	videoList := make([]*pb.Video, 0)
	for _, video := range videos {
		videoList = append(videoList, &pb.Video{
			Id:            video.Id,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			Title:         video.Title,
		})
	}
	resp := &pb.GetPublishVideoListResponse{
		VideoList: videoList,
	}
	return resp, nil
}
