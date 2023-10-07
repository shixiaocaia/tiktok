package dao

import (
	"github.com/shixiaocaia/tiktok/cmd/usersvr/dao/mysql"
	"gorm.io/gorm"
)

func UpdateWorkCount(authorID int64) error {
	db := mysql.GetDB()
	err := db.Model(&User{}).Where("id = ?", authorID).Update("work_count", gorm.Expr("work_count + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}
