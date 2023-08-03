package model

type User struct {
	// gorm.Model
	Id              int64  `gorm:"column:id; primary_key; auto_increment"`
	Name            string `gorm:"column:user_name"`
	Password        string `gorm:"column:password"`
	Follow          int64  `gorm:"column:follow_count"`
	Follower        int64  `gorm:"column:follower_count"`
	Avatar          string `gorm:"column:avatar"`
	BackgroundImage string `gorm:"column:background_image"`
	Signature       string `gorm:"column:signature"`
	TotalFav        int64  `gorm:"column:total_favorited"` // 我 获赞的总数
	FavCount        int64  `gorm:"column:favorite_count"`  // 我 点赞的视频总数
}

func (r *User) TableName() string {
	return "t_user"
}
