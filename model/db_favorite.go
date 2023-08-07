package model

type Favorite struct {
	Id      int64 `gorm:"column:id; primary_key;"` // favorite_id
	UserId  int64 `gorm:"column:user_id"`          // user_id 谁点的赞
	VideoId int64 `gorm:"column:video_id"`         // video_id 被点赞的视频
}

func (Favorite) TableName() string {
	return "t_favorite"
}
