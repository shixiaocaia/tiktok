package dao

type Message struct {
	Id         int64  `gorm:"column:id; primary_key;"`
	ToUserId   int64  `gorm:"column:to_user_id"`
	FromUserId int64  `gorm:"column:from_user_id"`
	Content    string `gorm:"column:content"`
	CreateTime int64  `gorm:"column:create_time"`
}

type NewMessage struct {
	Message string
	MsgType int64
}

func (r *Message) TableName() string {
	return "t_message"
}
