package utils

import (
	"context"
	"fmt"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/config"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"google.golang.org/grpc"
	"time"

	// 必须要导入这个包，否则grpc会报错
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
)

var (
	UserSvrClient pb.UserServiceClient
)

func GetCurrentTime() int64 {
	return time.Now().UnixNano() / 1e6
}

func InitSvrConn() {
	UserSvrClient = NewUserSvrClient(config.GetGlobalConfig().SvrConfig.UserSvrName)

}

func NewSvrConn(svrName string) (*grpc.ClientConn, error) {
	consulInfo := config.GetGlobalConfig().ConsulConfig
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) // 10秒超时
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, svrName),
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Errorf("NewSvrConn with svrname %s err:%v", svrName, err)
		return nil, err
	}
	log.Info("NewSvrConn success")
	return conn, nil
}

func NewUserSvrClient(svrName string) pb.UserServiceClient {
	conn, err := NewSvrConn(svrName)
	if err != nil {
		return nil
	}
	return pb.NewUserServiceClient(conn)
}