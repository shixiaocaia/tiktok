package dao

import (
	"gorm.io/gorm"
)

func UpdateFavorite(action, vid int64) error {
	db := GetDB()
	var num int64
	if action == 1 {
		num = 1
	} else {
		num = -1
	}
	err := db.Model(&Video{}).Where("id = ?", vid).Update("favorite_count", gorm.Expr("favorite_count + ?", num)).Error
	if err != nil {
		return err
	}
	return nil
}
