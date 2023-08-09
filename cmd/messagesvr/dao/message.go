package dao

import (
	"github.com/shixiaocaia/tiktok/cmd/messagesvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"gorm.io/gorm"
	"time"
)

func GetMessage(toUserId, fromUserId, preMsgTime int64) ([]*pb.Message, error) {
	db := GetDB()
	var messageList []*pb.Message
	log.Debugf("toUserId: %v, fromUserId: %v, preMsgTime: %v", toUserId, fromUserId, preMsgTime)
	err := db.Model(&Message{}).Where("((to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?)) AND create_time < ?", toUserId, fromUserId, fromUserId, toUserId, preMsgTime).Order("create_time").Limit(20).Find(&messageList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return messageList, nil
}

func InsertMessage(toUserId, fromUserId int64, content string) error {
	db := GetDB()
	message := Message{
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		Content:    content,
		CreateTime: time.Now().UnixNano() / 1e6,
	}
	err := db.Create(&message).Error
	if err != nil {
		return err
	}
	return nil
}

func GetNewestMessage(toUserId, fromUserId int64) (*NewMessage, error) {
	db := GetDB()
	var message Message
	err := db.Where("((to_user_id = ? AND from_user_id = ?) OR (to_user_id = ? AND from_user_id = ?))", toUserId, fromUserId, fromUserId, toUserId).Order("create_time desc").First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &NewMessage{
				Message: "",
				MsgType: 0,
			}, nil
		}
		return nil, err
	}
	newMessage := NewMessage{
		Message: message.Content,
	}
	if fromUserId == message.FromUserId {
		newMessage.MsgType = 1
	} else {
		newMessage.MsgType = 0
	}

	return &newMessage, nil
}
