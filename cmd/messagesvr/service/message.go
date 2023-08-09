package service

import (
	"context"
	"github.com/shixiaocaia/tiktok/cmd/messagesvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/messagesvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

type MessageService struct {
	pb.UnimplementedMessageServiceServer
}

func (u MessageService) MessageChat(ctx context.Context, req *pb.MessageChatReq) (*pb.MessageChatRsp, error) {
	log.Debugf("MessageChat req: %v", req)
	messageList, err := dao.GetMessage(req.ToUserId, req.FromUserId, req.PreMsgTime)
	if err != nil {
		log.Errorf("GetMessage failed: %v", err)
		return nil, err
	}

	return &pb.MessageChatRsp{
		MessageList: messageList,
	}, nil

}

func (u MessageService) MessageAction(ctx context.Context, req *pb.MessageActionReq) (*pb.MessageActionRsp, error) {
	err := dao.InsertMessage(req.ToUserId, req.FromUserId, req.Content)
	if err != nil {
		log.Errorf("InsertMessage failed: %v", err)
		return nil, err
	}
	return &pb.MessageActionRsp{}, nil
}

func (u MessageService) NewestMessageDic(ctx context.Context, req *pb.NewestMessageReq) (*pb.NewestMessageRsp, error) {
	var rsp = &pb.NewestMessageRsp{
		NewestMessageDict: make(map[int64]*pb.NewestMessage),
	}

	for _, friend := range req.FriendIdList {
		newestMessage, err := dao.GetNewestMessage(req.UserId, friend)
		if err != nil {
			log.Errorf("GetNewestMessage failed: %v", err)
			return nil, err
		}
		rsp.NewestMessageDict[friend] = &pb.NewestMessage{
			Message: newestMessage.Message,
			MsgType: newestMessage.MsgType,
		}
	}
	return rsp, nil
}
