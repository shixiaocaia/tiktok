package main

import (
	"fmt"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/config"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/routes"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/utils"
	"go.uber.org/zap"
)

func Init() {
	// 读取配置
	if err := config.Init(); err != nil {
		log.Fatalf("init gatewaysvr config failed, err:%v\n", err)
	}
	// 初始化日志
	log.InitLogger()
	//log.Test("www.baidu.com")
	log.Info("log init success")
	// 初始化微服务
	utils.InitSvrConn()
}

func main() {
	Init()
	defer log.Sync()

	//初始化路由
	r := routes.SetRoute()
	//启动
	if err := r.Run(fmt.Sprintf(":%d", config.GetGlobalConfig().SvrConfig.Port)); err != nil {
		zap.L().Panic("Router.Run error: ", zap.Error(err))
	}

}
