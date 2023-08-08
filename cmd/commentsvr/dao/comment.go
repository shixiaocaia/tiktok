package dao

import (
	"gorm.io/gorm"
	"time"
)

func CommentAdd(uid, vid int64, text string) (*Comment, error) {
	db := GetDB()
	nowTime := time.Now()
	comment := Comment{
		UserId:      uid,
		VideoId:     vid,
		CommentText: text,
		CreateTime:  nowTime,
	}

	err := db.Create(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
	// todo redis缓存
}

func CommentDel(commentId, vid int64) error {
	db := GetDB()
	err := db.Where("id = ? and video_id = ?", commentId, vid).Delete(&Comment{}).Error
	if err != nil {
		return err
	}

	// todo redis缓存更新

	return nil
}

func GetCommentList(vid int64) ([]*Comment, error) {
	db := GetDB()
	var commentList []*Comment
	err := db.Model(&Comment{}).Where("video_id = ?", vid).Find(&commentList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*Comment{}, nil
		}
		return nil, err
	}
	return commentList, nil

}
