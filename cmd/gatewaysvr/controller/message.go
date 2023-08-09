package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/constant"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/response"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"strconv"
)

func MessageChat(ctx *gin.Context) {
	UserId, _ := ctx.Get("UserID")
	toUserId, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
	if err != nil {
		log.Errorf("to_user_id is not int64: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	lastTime, err := strconv.ParseInt(ctx.Query("pre_msg_time"), 10, 64)
	if err != nil || lastTime == int64(0) {
		log.Errorf("pre_msg_time is invalid")
		lastTime = utils.GetCurrentTime()
	}
	log.Debugf("userId: %v, toUserId: %v, pre_msg_time: %v", UserId, toUserId, lastTime)

	messageChatRsp, err := utils.GetMessageSvrClient().MessageChat(ctx, &pb.MessageChatReq{
		ToUserId:   toUserId,
		FromUserId: UserId.(int64),
		PreMsgTime: lastTime,
	})
	//log.Debugf("messageChatRsp: %v", messageChatRsp)
	if err != nil {
		log.Errorf("MessageChat failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	//rsp := &pb.MessageChatRsp{
	//	MessageList: make([]*pb.Message, 0),
	//}
	//for _, message := range messageChatRsp.MessageList {
	//	rsp.MessageList = append(rsp.MessageList, &pb.Message{
	//		Id:         message.Id,
	//		ToUserId:   message.ToUserId,
	//		FromUserId: message.FromUserId,
	//		Content:    message.Content,
	//		CreateTime: strconv.FormatInt(message.CreateTime.UnixNano()/1e6, 10),
	//	})
	//}

	log.Infof("Get MessageChat success...")
	response.Success(ctx, "success", messageChatRsp)

}

func MessageAction(ctx *gin.Context) {
	var message constant.Message
	if err := ctx.ShouldBind(&message); err != nil {
		log.Errorf("ShouldBind(&message) failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	message.FromUserId = ctx.GetInt64("UserID")

	rsp, err := utils.GetMessageSvrClient().MessageAction(ctx, &pb.MessageActionReq{
		ToUserId:   message.ToUserId,
		FromUserId: message.FromUserId,
		Content:    message.Content,
	})

	if err != nil {
		log.Errorf("MessageAction failed: %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	log.Infof("MessageAction success...")
	response.Success(ctx, "success", rsp)
}
