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

func IsFavoriteVideo(uid, vid int64) (bool, error) {
	db := GetDB()
	err := db.Where("user_id = ? and video_id = ?", uid, vid).First(&model.Favorite{}).Error
	if err != nil {
		// 没有找到记录说明不是点赞视频
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetFavoriteVideoIdList(uid int64) ([]int64, error) {
	db := GetDB()
	var favoriteList []*model.Favorite
	err := db.Model(&model.Favorite{}).Where("user_id = ?", uid).Find(&favoriteList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []int64{}, nil
		}
		return nil, err
	}

	videoIdList := make([]int64, 0)
	for _, video := range favoriteList {
		videoIdList = append(videoIdList, video.VideoId)
	}

	return videoIdList, nil
}
