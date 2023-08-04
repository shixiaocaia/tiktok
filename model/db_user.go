package model

type User struct {
	// gorm.Model
	Id              int64  `gorm:"column:id; primary_key; auto_increment"` // 用户id
	Name            string `gorm:"column:user_name"`                       // 用户名
	Password        string `gorm:"column:password"`                        // 密码
	Follow          int64  `gorm:"column:follow_count"`                    // 关注数
	Follower        int64  `gorm:"column:follower_count"`                  // 粉丝数
	Avatar          string `gorm:"column:avatar"`                          // 头像
	BackgroundImage string `gorm:"column:background_image"`                // 背景图
	Signature       string `gorm:"column:signature"`                       // 签名
	TotalFav        int64  `gorm:"column:total_favorited"`                 // 获赞的总数
	FavCount        int64  `gorm:"column:favorite_count"`                  // 点赞的视频总数
	WorkCount       int64  `gorm:"column:work_count"`                      // 发布的视频总数
}

func (r *User) TableName() string {
	return "t_user"
}
