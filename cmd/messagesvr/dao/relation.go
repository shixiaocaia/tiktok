package dao

import (
	"fmt"
	"github.com/shixiaocaia/tiktok/cmd/relationsvr/constant"
	"github.com/shixiaocaia/tiktok/cmd/relationsvr/log"
	"gorm.io/gorm"
)

func FollowAdd(follow, follower int64) error {
	db := GetDB()
	relation := Relation{
		Follow:   follow,
		Follower: follower,
	}

	// 1. 先查询是否已关注，避免重复关注
	if err := db.Where("follow_id = ? and follower_id = ?", follow, follower).First(&Relation{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 记录不存在，插入
			err = db.Create(&relation).Error
			if err != nil {
				log.Errorf("create failed: %v", err)
				return err
			}
		} else {
			return err
		}
	} else {
		return fmt.Errorf(constant.HavaFollowed)
	}
	return nil
}

func FollowDel(follow, follower int64) error {
	db := GetDB()

	err := db.Where("follow_id = ? and follower_id = ?", follow, follower).Delete(&Relation{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("Delete failed: %v", err)
		return err
	}
	return nil
}

func GetFollowIdList(uid int64) ([]*Relation, error) {
	db := GetDB()
	relationList := make([]*Relation, 0)
	err := db.Where("follower_id = ?", uid).Find(&relationList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return relationList, nil
}

func GetFollowerIdList(uid int64) ([]*Relation, error) {
	db := GetDB()
	relationList := make([]*Relation, 0)
	err := db.Where("follow_id = ?", uid).Find(&relationList).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return relationList, nil
}

func IsFriend(followId, followerId int64) (bool, error) {
	db := GetDB()
	err := db.Where("follow_id = ? and follower_id = ?", followId, followerId).First(&Relation{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
