package dao

import (
	"github.com/shixiaocaia/tiktok/cmd/videosvr/log"
	"github.com/shixiaocaia/tiktok/model"
	"gorm.io/gorm"
)

func GetVideoListByFeed(time int64) ([]model.Video, error) {
	var videos []model.Video
	db := GetDB()
	err := db.Where("publish_time < ?", time).Limit(20).Order("publish_time DESC").Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("can not find videos...")
		return videos, err
	}
	log.Info("GetVideoListByFeed...")
	return videos, nil
}
