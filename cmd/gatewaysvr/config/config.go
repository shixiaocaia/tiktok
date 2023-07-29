package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/shixiaocaia/tiktok/cmd/gatewaysvr/log"
	"github.com/spf13/viper"
)

var globalConfig = new(GlobalConfig)

type GlobalConfig struct {
	Port int
	Ping string
}

// Init 初始化配置
func Init() (err error) {
	viper.SetConfigFile("./config/config.yaml") // 指定配置文件路径
	err = viper.ReadInConfig()                  // 读取配置信息
	if err != nil {                             // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(globalConfig); err != nil {
		log.Error("viper.ReadInConfig() failed")
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}
	// 监控配置文件变化
	viper.WatchConfig()
	// 注意！！！配置文件发生变化后要同步到全局变量Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Info("配置信息更新...")
		if err := viper.Unmarshal(globalConfig); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}
	})
	return
}

func GetGlobalConfig() *GlobalConfig {
	return globalConfig
}
