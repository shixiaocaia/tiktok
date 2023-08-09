package dao

import (
	"github.com/shixiaocaia/tiktok/cmd/videosvr/log"
	"gorm.io/gorm"
	"time"
)

// GetVideoListByFeed 获取视频信息
func GetVideoListByFeed(time int64) ([]Video, error) {
	var videos []Video
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
	video := Video{
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
func GetVideoListByAuthorID(authorId int64) ([]Video, error) {
	var videos []Video

	db := GetDB()
	err := db.Where("author_id = ?", authorId).Order("id DESC").Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("GetVideoListByAuthorID failed: %v", err)
		return nil, err
	}
	log.Debugf("videos: %v", videos)
	return videos, nil
}

// GetVideoListByVideoIdList 获取用户发布的多个视频
func GetVideoListByVideoIdList(videoId []int64) ([]Video, error) {
	var videos []Video
	db := GetDB()
	err := db.Where("id in ?", videoId).Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return videos, err
	}
	return videos, nil
}

// UpdateCommentCount 更新评论数
func UpdateCommentCount(vid, actionType int64) error {
	num := 1
	if actionType == 2 {
		num = -1
	}

	db := GetDB()
	err := db.Model(&Video{}).Where("id = ?", vid).Update("comment_count", gorm.Expr("comment_count + ?", num)).Error
	if err != nil {
		return err
	}
	return nil
}
