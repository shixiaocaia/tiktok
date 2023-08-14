package dao

type Favorite struct {
	Id      int64 `gorm:"column:id; primary_key;"` // favorite_id
	UserId  int64 `gorm:"column:user_id"`          // user_id 谁点的赞
	VideoId int64 `gorm:"column:video_id"`         // video_id 被点赞的视频
}

func (Favorite) TableName() string {
	return "t_favorite"
}

type Video struct {
	Id            int64  `gorm:"column:id; primary_key;"` // video_id
	AuthorId      int64  `gorm:"column:author_id;"`       // 谁发布的
	PlayUrl       string `gorm:"column:play_url;"`        // videoURL
	CoverUrl      string `gorm:"column:cover_url;"`       // picURL
	FavoriteCount int64  `gorm:"column:favorite_count;"`  // 点赞数
	CommentCount  int64  `gorm:"column:comment_count;"`   // 评论数
	PublishTime   int64  `gorm:"column:publish_time;"`    // 发布时间
	Title         string `gorm:"column:title;"`           // 标题
}

func (Video) TableName() string {
	return "t_video"
}
