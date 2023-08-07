package dao

import (
	"github.com/shixiaocaia/tiktok/cmd/videosvr/log"
	"github.com/shixiaocaia/tiktok/model"
	"gorm.io/gorm"
	"time"
)

// GetVideoListByFeed 获取视频信息
func GetVideoListByFeed(time int64) ([]model.Video, error) {
	var videos []model.Video
	db := GetDB()
	err := db.Where("publish_time < ?", time).Limit(30).Order("publish_time DESC").Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("can not find videos...")
		return videos, err
	}
	log.Info("GetVideoListByFeed...")
	return videos, nil
}

// InsertVideo 记录视频信息
func InsertVideo(authorID int64, playUrl, picUrl, title string) error {
	video := model.Video{
		AuthorId:      authorID,
		PlayUrl:       playUrl,
		CoverUrl:      picUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		PublishTime:   time.Now().UnixNano() / 1e6,
		Title:         title,
	}
	db := GetDB()
	err := db.Create(&video).Error
	if err != nil {
		log.Errorf("db.Create failed: %v", err)
		return err
	}
	return nil

}

// GetVideoListByAuthorID 根据authorID获取视频
func GetVideoListByAuthorID(authorId int64) ([]model.Video, error) {
	var videos []model.Video

	db := GetDB()
	err := db.Where("author_id = ?", authorId).Order("id DESC").Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("GetVideoListByAuthorID faild: %v", err)
		return nil, err
	}
	log.Debugf("videos: %v", videos)
	return videos, nil
}

// GetVideoListByVideoIdList 获取用户发布的多个视频
func GetVideoListByVideoIdList(videoId []int64) ([]model.Video, error) {
	var videos []model.Video
	db := GetDB()
	err := db.Where("id in ?", videoId).Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return videos, err
	}
	return videos, nil
}
