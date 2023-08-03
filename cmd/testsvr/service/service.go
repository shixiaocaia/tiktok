package service

import (
	"context"
	"fmt"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"log"
)

// grpc的具体逻辑

type GreetService struct {
	pb.UnimplementedGreeterServer
}

// // UpdateUserFavoritedCount 更新用户 获赞数，ActionType 1：表示+1 2：-1
//
//	func (u UserService) UpdateUserFavoritedCount(ctx context.Context, req *pb.UpdateUserFavoritedCountReq) (*pb.UpdateUserFavoritedCountRsp, error) {
//		err := repository.UpdateUserFavoritedNum(req.UserId, req.ActionType)
//		if err != nil {
//			log.Errorf("UpdateUserFavoritedCount err", req.UserId)
//			return nil, err
//		}
//		return &pb.UpdateUserFavoritedCountRsp{}, nil
//	}
func (*GreetService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	msg := fmt.Sprintf("I got message from request: %s\n", req.Name)
	log.Println(msg)
	return &pb.HelloReply{Message: msg}, nil
}
