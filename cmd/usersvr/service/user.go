package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/constant"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/dao"
	"github.com/shixiaocaia/tiktok/cmd/usersvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"golang.org/x/crypto/bcrypt"
)

// JWT 签名密钥，定义JWT中存放usrname和userID

var mySigningKey = []byte("mini_tiktok")

type JWTClaims struct {
	Username string `json:"user_name"`
	UserID   int64  `json:"user_id"`
	jwt.RegisteredClaims
}

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
	token, err := GenToken(info.Id, req.Username)
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
	token, err := GenToken(info.Id, info.Name)
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

func GenToken(userid int64, username string) (string, error) {
	claims := &JWTClaims{
		Username: username,
		UserID:   userid,
		RegisteredClaims: jwt.RegisteredClaims{
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			// 方便测试，不设置过期时间
			Issuer: "server",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
