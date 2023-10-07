package dao

import (
	"github.com/shixiaocaia/tiktok/cmd/usersvr/dao/mysql"
	"gorm.io/gorm"
)

// UpdateFavorited 更新被喜欢数
func UpdateFavorited(uid, action int64) error {
	db := mysql.GetDB()
	var num int64
	// updateType 1: 点赞 else： 取消点赞
	if action == 1 {
		num = 1
	} else {
		num = -1
	}
	err := db.Model(&User{}).Where("id = ?", uid).Update("total_favorited", gorm.Expr("total_favorited + ?", num)).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateFavorite 更新自己喜欢的作品数
func UpdateFavorite(uid, action int64) error {
	db := mysql.GetDB()
	var num int64
	// updateType 1: 点赞 else： 取消点赞
	if action == 1 {
		num = 1
	} else {
		num = -1
	}
	err := db.Model(&User{}).Where("id = ?", uid).Update("favorite_count", gorm.Expr("favorite_count + ?", num)).Error
	if err != nil {
		return err
	}
	return nil
}
