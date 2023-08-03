package service

import (
	"context"
	"fmt"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
)

// grpc的具体逻辑

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (u UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 检查用户是否已经存在
	sign, err := dao.UserNameIsExist(req.Username)
	// 函数执行错误
	if err != nil {
		log.Errorf("UserNameIsExist failed: %v", err)
		return nil, err
	}
	// 用户名已经存在
	if sign {
		return nil, fmt.Errorf("username %s exists", req.Username)
	}
	// 写入mysql
	info, err := dao.InsertUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	// 生成token

	resp := &pb.RegisterResponse{
		UserId: info.Id,
		Token:  "2233",
	}
	return resp, nil
}
