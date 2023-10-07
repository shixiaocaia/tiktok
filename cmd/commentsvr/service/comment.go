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
		log.Infof("user:%v comment: %v", comment.UserId, comment.CommentText)
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

	// 更新缓存
	if err := dao.DelCommentCacheInfo(req.VideoId); err != nil {
		log.Errorf("DelCommentCacheInfo failed: %v", err)
		return nil, err
	}
	return rsp, nil
}

// GetCommentList 获取评论列表
func (u CommentService) GetCommentList(ctx context.Context, req *pb.GetCommentListReq) (*pb.GetCommentListRsp, error) {
	videoId := req.VideoId
	// 1. 从缓存中获取评论列表
	list, err := dao.GetCommentCacheList(videoId)
	if len(list) == 0 || err != nil {
		// 2. 缓存未命中，从数据库中获取评论列表
		log.Infof("GetCommentCacheList failed: %v", err)

		list, err = dao.GetCommentList(videoId)
		if err != nil {
			log.Errorf("GetCommentList failed: %v", err)
			return nil, err
		}
		// 3. 将评论列表写入缓存
		if len(list) != 0 {
			err = dao.SetCommentCacheInfo(videoId, list)
			if err != nil {
				log.Errorf("SetCommentCacheInfo failed: %v", err)
				return nil, err
			}
			log.Info("SetCommentCacheInfo success")
		}
	} else {
		log.Info("GetCommentCacheList success")
	}
	// 3. 构建返回值
	var rsp = &pb.GetCommentListRsp{}
	for _, comment := range list {
		rsp.CommentList = append(rsp.CommentList, BuildComment(comment))
	}

	return rsp, err
}

// BuildComment 构建评论信息
func BuildComment(comment *dao.Comment) *pb.Comment {
	return &pb.Comment{
		Id:         comment.Id,
		Content:    comment.CommentText,
		CreateDate: comment.CreateTime.Format(constant.DefaultTime),
		User: &pb.UserInfo{
			Id: comment.UserId,
		},
	}
}
