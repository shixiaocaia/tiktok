package service

import (
	"context"
	"github.com/shixiaocaia/tiktok/cmd/commentsvr/constant"
	"github.com/shixiaocaia/tiktok/cmd/commentsvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/commentsvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer
}

func (u CommentService) CommentAction(ctx context.Context, req *pb.CommentActionReq) (*pb.CommentActionRsp, error) {
	// 评论 or 删评
	var rsp = &pb.CommentActionRsp{
		Comment: nil,
	}
	if req.ActionType == 1 {
		comment, err := dao.CommentAdd(req.UserId, req.VideoId, req.CommentText)
		if err != nil {
			log.Errorf("CommentAdd failed: %v", err)
			return nil, err
		}
		log.Debugf("comment: %v", comment)
		rsp.Comment = &pb.Comment{
			Id:         comment.Id,
			User:       nil,
			Content:    comment.CommentText,
			CreateDate: comment.CreateTime.Format(constant.DefaultTime),
		}
	} else {
		err := dao.CommentDel(req.CommentId, req.VideoId)
		if err != nil {
			log.Errorf("CommentDel failed: %v", err)
			return nil, err
		}
	}
	return rsp, nil
}

func (u CommentService) GetCommentList(ctx context.Context, req *pb.GetCommentListReq) (*pb.GetCommentListRsp, error) {
	videoId := req.VideoId
	list, err := dao.GetCommentList(videoId)
	if err != nil {
		log.Errorf("GetCommentList failed: %v", err)
		return nil, err
	}

	var rsp = &pb.GetCommentListRsp{}
	for _, comment := range list {
		rsp.CommentList = append(rsp.CommentList, &pb.Comment{
			Id:         comment.Id,
			Content:    comment.CommentText,
			CreateDate: comment.CreateTime.Format(constant.DefaultTime),
			User: &pb.UserInfo{
				Id: comment.UserId,
			},
		})
	}

	return rsp, err
}
