package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/constant"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/log"
	jwtoken "github.com/shixiaocaia/tiktok/cmd/usersvr/middleware/jwt"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"golang.org/x/crypto/bcrypt"
)

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
		log.Error("UserNameIsExist ", req.Username)
		return nil, fmt.Errorf(constant.UserNameIsExist)
	}
	// 写入mysql
	info, err := dao.InsertUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	// 生成token
	token, err := jwtoken.GenToken(info.Id, req.Username)
	if err != nil {
		return nil, err
	}

	resp := &pb.RegisterResponse{
		UserId: info.Id,
		Token:  token,
	}
	return resp, nil
}

func (u UserService) CheckPassWord(ctx context.Context, req *pb.CheckPassWordRequest) (*pb.CheckPassWordResponse, error) {
	// 查看用户是否存在
	info, err := dao.GetUserInfo(req.Username)
	if err != nil {
		log.Error(constant.ErrorUserInfo)
		return nil, err
	}
	// 用户存在查看密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(info.Password), []byte(req.Password))
	if err != nil {
		log.Errorf(constant.ErrorPassword)
		return nil, errors.New(constant.ErrorPassword)
	}
	// 生成token,后续验证
	token, err := jwtoken.GenToken(info.Id, info.Name)
	if err != nil {

		return nil, err
	}
	response := &pb.CheckPassWordResponse{
		UserId: info.Id,
		Token:  token,
	}
	log.Info("login success...")
	return response, nil

}

func (u UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	info, err := dao.GetUserInfo(req.UserId)
	if err != nil {
		log.Error(constant.ErrorUserInfo)
		return nil, err
	}
	response := &pb.GetUserInfoResponse{
		User: UserToUserInfo(info),
	}
	log.Info("getUserinfo success...")
	return response, nil
}

func (u UserService) GetUserInfoDict(ctx context.Context, req *pb.GetUserInfoDictRequest) (*pb.GetUserInfoDictResponse, error) {
	userIdList := req.UserIdList
	log.Debugf("userIdList ", userIdList)
	// 应该获取用户ID，一次性查询而不是单个查询，效率低
	userList, err := dao.GetUserListInfo(userIdList)
	if err != nil {
		log.Error("GetUserListInfo failed...")
		return nil, err
	}
	log.Debugf("userList: %v", userList)
	response := &pb.GetUserInfoDictResponse{UserInfoDict: make(map[int64]*pb.UserInfo)}

	for _, user := range userList {
		response.UserInfoDict[user.Id] = UserToUserInfo2(user)
	}
	log.Info("GetUserInfoDict success...")
	return response, nil

}

// UserToUserInfo 类型转换
func UserToUserInfo(info dao.User) *pb.UserInfo {
	return &pb.UserInfo{
		Id:              info.Id,
		Name:            info.Name,
		FollowCount:     info.Follow,
		FollowerCount:   info.Follower,
		IsFollow:        false,
		Avatar:          info.Avatar,
		BackgroundImage: info.BackgroundImage,
		Signature:       info.Signature,
		TotalFavorited:  info.TotalFav,
		FavoriteCount:   info.FavCount,
		WorkCount:       info.WorkCount,
	}
}

// UserToUserInfo2 类型转换(区别上面指针）
func UserToUserInfo2(info *dao.User) *pb.UserInfo {
	return &pb.UserInfo{
		Id:              info.Id,
		Name:            info.Name,
		FollowCount:     info.Follow,
		FollowerCount:   info.Follower,
		IsFollow:        false,
		Avatar:          info.Avatar,
		BackgroundImage: info.BackgroundImage,
		Signature:       info.Signature,
		TotalFavorited:  info.TotalFav,
		FavoriteCount:   info.FavCount,
		WorkCount:       info.WorkCount,
	}
}
