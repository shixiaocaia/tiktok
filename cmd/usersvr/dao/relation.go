package dao

import (
	"github.com/shixiaocaia/tiktok/cmd/usersvr/log"
	"gorm.io/gorm"
)

func UpdateFollowCount(uid, actionType int64) error {
	num := 1
	if actionType == 2 {
		num = -1
	}
	log.Debugf("UpdateFollowCount")
	db := GetDB()
	err := db.Model(&User{}).Where("id = ?", uid).Update("follow_count", gorm.Expr("follow_count + ?", num)).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateFollowerCount(uid, actionType int64) error {
	num := 1
	if actionType == 2 {
		num = -1
	}

	db := GetDB()
	err := db.Model(&User{}).Where("id = ?", uid).Update("follower_count", gorm.Expr("follower_count + ?", num)).Error
	if err != nil {
		return err
	}
	return nil
}
