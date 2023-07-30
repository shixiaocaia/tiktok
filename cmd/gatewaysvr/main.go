package main

import (
	"fmt"
	"gatewaysvr/config"
	"gatewaysvr/log"
	"gatewaysvr/routes"
	"go.uber.org/zap"
)

func Init() {
	// 读取配置
	if err := config.Init(); err != nil {

	}
	// 初始化日志
	log.InitLogger()
	//log.Test("www.baidu.com")
	log.Info("log init success")

	// 初始化微服务
}

func main() {
	Init()

	// 初始化路由
	r := routes.SetRoute()
	// 启动
	if err := r.Run(fmt.Sprintf(":%d", config.GetGlobalConfig().Port)); err != nil {
		zap.L().Panic("Router.Run error: ", zap.Error(err))
	}

}
