package main

import (
	"fmt"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/config"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/routes"
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
	r.Run(fmt.Sprintf(":%d", config.GetGlobalConfig().Port))

}
