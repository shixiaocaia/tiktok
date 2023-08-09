package dao

import (
	"gorm.io/gorm"
)

func UpdateWorkCount(authorID int64) error {
	db := GetDB()
	err := db.Model(&User{}).Where("id = ?", authorID).Update("work_count", gorm.Expr("work_count + ?", 1)).Error
	if err != nil {
		return err
	}
	return nil
}
