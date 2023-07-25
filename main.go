package main

import (
	"go.uber.org/zap"
	"main/config"
	"main/log"
	"main/router"
)

func init() {
	// 1. 读取参数
	if err := config.Init(); err != nil {
		log.Fatalf("config.Init() failed: %v\n", zap.Error(err))
		return
	}

	// 2. 初始化日志
	err := log.InitLog()
	if err != nil {
		log.Fatalf("log.InitLog() failed: %v\n", zap.Error(err))
	}
	defer log.Sync()
}

func main() {

	// 3. 初始化路由
	r := router.SetRouter()
	// 启动服务， spin支持优雅退出
	r.Spin()

}
