package dao

import (
	"github.com/shixiaocaia/tiktok/cmd/favoritesvr/log"
	"github.com/shixiaocaia/tiktok/model"
	"gorm.io/gorm"
)

func LikeAction(userID, videoID int64) error {
	db := GetDB()
	favorite := &model.Favorite{
		UserId:  userID,
		VideoId: videoID,
	}
	err := db.Where("user_id = ? and video_id = ?", userID, videoID).First(&model.Favorite{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	// 确保没有记录，插入新记录
	err = db.Create(&favorite).Error
	log.Debugf("LikeAction: %v", favorite)
	if err != nil {
		return err
	}
	return nil
}

func DislikeAction(userID, videoID int64) error {
	db := GetDB()
	err := db.Where("user_id = ? and video_id = ?", userID, videoID).Delete(&model.Favorite{}).Error
	if err != nil {
		return err
	}
	log.Infof("Dislike video: %v", videoID)
	return nil
}
