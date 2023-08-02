package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/shixiaocaia/tiktok/cmd/testsvr/config"
	"github.com/shixiaocaia/tiktok/cmd/testsvr/log"
	"github.com/shixiaocaia/tiktok/cmd/testsvr/service"
	"github.com/shixiaocaia/tiktok/cmd/testsvr/utils/consul"
	"github.com/shixiaocaia/tiktok/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	_ "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Init() {
	// 读取配置
	if err := config.Init(); err != nil {

		log.Fatalf("init config failed, err:%v\n", err)
	}
	// 初始化日志
	log.InitLogger()

	log.Info("log init success")

}

func Run() error {
	//log.Debugf(config.GetGlobalConfig().SvrConfig.Host, config.GetGlobalConfig().SvrConfig.Port)
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "", config.GetGlobalConfig().SvrConfig.Port))
	if err != nil {
		log.Fatalf("listen error: %v", err)
		return fmt.Errorf("listen error: %v", err)
	}

	// 启动grpc server
	server := grpc.NewServer()
	// 注册grpc server
	pb.RegisterGreeterServer(server, &service.GreetService{})
	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 注册服务到consul中
	consulClient := consul.NewRegistryClient(config.GetGlobalConfig().ConsulConfig.Host, config.GetGlobalConfig().ConsulConfig.Port)
	// 服务ID指定而不是生成
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	if err := consulClient.Register(config.GetGlobalConfig().SvrConfig.Host, config.GetGlobalConfig().SvrConfig.Port,
		config.GetGlobalConfig().Name, config.GetGlobalConfig().ConsulConfig.Tags, serviceID); err != nil {
		log.Fatal("consul.Register error: ", err)
		return fmt.Errorf("consul.Register error: ", err)
	}
	log.Info("Init Consul Register success")

	// 启动
	log.Infof("TikTokLite.usersvr listening on %s:%d", config.GetGlobalConfig().SvrConfig.Host, config.GetGlobalConfig().SvrConfig.Port)
	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 服务终止，注销 consul 服务
	if err = consulClient.DeRegister(serviceID); err != nil {
		log.Info("注销失败")
		return fmt.Errorf("注销失败")
	} else {
		log.Info("注销成功")
	}
	return nil
}

func main() {
	Init()
	defer log.Sync()

	//log.Debugf("%v", config.GetGlobalConfig().SvrConfig.Port)
	if err := Run(); err != nil {
		log.Errorf("Usersvr run err: %v", err)
	}
}
