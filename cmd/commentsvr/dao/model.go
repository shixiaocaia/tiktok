package dao

import "time"

type Comment struct {
	Id          int64     `gorm:"column:id; primary_key;"` // ID = commentID
	UserId      int64     `gorm:"column:user_id"`          // 是谁评论的
	VideoId     int64     `gorm:"column:video_id"`         // 评论哪个视频
	CommentText string    `gorm:"column:comment_text"`     // 评论内容
	CreateTime  time.Time `gorm:"column:create_time"`      // 评论时间
}

func (c *Comment) TableName() string {
	return "t_comment"
}
