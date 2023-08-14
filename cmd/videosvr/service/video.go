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

const (
	maxRetries = 3
	retryDelay = 500 * time.Millisecond
)

// GetFeedList 获取一组视频基本信息
func (u VideoService) GetFeedList(ctx context.Context, req *pb.GetFeedListRequest) (*pb.GetFeedListResponse, error) {
	// 返回以时间戳为止的一组视频
	// 1. 查缓存
	videoList := make([]dao.Video, 0)
	var err error
	videoList, err = dao.GetVideoListCache(req.LatestTime)
	if err != nil || len(videoList) == 0 {
		log.Info("GetVideoCache failed")
		videoList, err = dao.GetVideoListByFeed(req.LatestTime)
		if err != nil {
			log.Error("GetVideoListByFeed failed")
			return nil, err
		}
		// 更新缓存
		err = dao.SetVideoCache(videoList)
		if err != nil {
			log.Error("SetVideoCache failed")
		}
	} else {
		log.Info("GetVideoCache success")
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
	startTime := time.Now()
	client := minio.GetMinio()
	playUrl, err := client.UploadFile("video", req.SaveFile, strconv.FormatInt(req.UserId, 10))
	if err != nil {
		log.Errorf("Minio UploadFile err: %v", err)
		return nil, err
	}
	log.Infof("save file: %v, cost %v", req.SaveFile, time.Since(startTime))

	// 生成视频封面
	imageFile, err := utils.GetImageFile(req.SaveFile)
	if err != nil {
		log.Errorf("ffmpeg getVideoPic failed: %v", err)
		return nil, err
	}
	log.Infof("GetImageFile cost %v", time.Since(startTime))

	// 上传封面
	coverUrl, err := client.UploadFile("pic", imageFile, strconv.FormatInt(req.UserId, 10))
	if err != nil {
		log.Errorf("minio upLoadPic failed: %v", err)
		return nil, err
	}
	log.Infof("UpImageFile cost %v", time.Since(startTime))

	err = dao.InsertVideo(req.UserId, playUrl, coverUrl, req.Title)
	if err != nil {
		log.Errorf("InsertVideo failed: %v", err)
		return nil, err
	}

	// 更新缓存
	err = dao.UpdateVideoCache()
	if err != nil {
		log.Errorf("UpdateVideoCache failed: %v", err)
		return nil, err
	}

	return &pb.PublishVideoResponse{}, nil
}

// GetPublishVideoList 获得上传视频列表
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

// GetVideoInfoList 获取视频作者信息
func (u VideoService) GetVideoInfoList(ctx context.Context, req *pb.GetVideoInfoListReq) (*pb.GetVideoInfoListRsp, error) {
	videoList, err := dao.GetVideoListByVideoIdList(req.VideoId)
	if err != nil {
		log.Errorf("GetVideoListByVideoIdList failed...")
		return nil, err
	}

	rsp := &pb.GetVideoInfoListRsp{
		VideoInfoList: make([]*pb.VideoInfo, 0),
	}
	for _, video := range videoList {
		rsp.VideoInfoList = append(rsp.VideoInfoList, &pb.VideoInfo{
			Id:            video.Id,
			AuthorId:      video.AuthorId,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    false,
			Title:         video.Title,
		})
	}
	return rsp, nil
}

// UpdateCommentCount 更新视频评论数
func (u VideoService) UpdateCommentCount(ctx context.Context, req *pb.UpdateCommentCountReq) (*pb.UpdateCommentCountRsp, error) {
	err := dao.UpdateCommentCount(req.VideoId, req.ActionType)
	if err != nil {
		return nil, err
	}

	// 更新VideoInfo中评论数缓存
	for i := 0; i < maxRetries; i++ {
		err = dao.UpdateCommentCountCache(req.VideoId, req.ActionType)
		if err != nil {
			log.Errorf("UpdateCommentCacheInfo failed: %v", err)
			time.Sleep(retryDelay)
			continue
		}
		if err == nil {
			log.Info("UpdateVideoCache success")
			break
		}
	}
	// 更新缓存失败，清空缓存，改为旁路缓存
	if err != nil {
		err = dao.UpdateVideoCache()
		if err != nil {
			log.Errorf("del videoCache failed: %v", err)
			return nil, err
		}
	}

	return &pb.UpdateCommentCountRsp{}, nil
}

// UpdateFavoriteCount 更新视频点赞数
func (u VideoService) UpdateFavoriteCount(ctx context.Context, req *pb.UpdateFavoriteCountReq) (*pb.UpdateFavoriteCountRsp, error) {
	err := dao.UpdateFavorite(req.ActionType, req.VideoId)
	if err != nil {
		log.Errorf("UpdateFavoriteCount failed: %v", err)
		return nil, err
	}

	// 更新VideoInfo中点赞数缓存
	for i := 0; i < maxRetries; i++ {
		err = dao.UpdateFavoriteCountCache(req.VideoId, req.ActionType)
		if err != nil {
			log.Errorf("UpdateFavoriteCountCache failed: %v", err)
			time.Sleep(retryDelay)
			continue
		}
		if err == nil {
			log.Info("UpdateVideoCache success")
			break
		}
	}
	// 更新缓存失败，清空缓存，改为旁路缓存
	if err != nil {
		err = dao.UpdateVideoCache()
		if err != nil {
			log.Errorf("del videoCache failed: %v", err)
			return nil, err
		}
	}
	return &pb.UpdateFavoriteCountRsp{}, nil
}
