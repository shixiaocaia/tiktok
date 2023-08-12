package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var globalConfig = new(GlobalConfig)

type GlobalConfig struct {
	*SvrConfig    `mapstructure:"svr_config"`
	*ConsulConfig `mapstructure:"consul"`
	*LogConfig    `mapstructure:"log"`
	*MySQLConfig  `mapstructure:"mysql"`
	*RedisConfig  `mapstructure:"redis"`
}

type SvrConfig struct {
	Name string `mapstructure:"name"` // 服务name
	Host string `mapstructure:"host"` // 服务host
	Port int    `mapstructure:"port"` // 服务port
}

type ConsulConfig struct {
	Host string   `mapstructure:"host"`
	Port int      `mapstructure:"port"`
	Tags []string `mapstructure:"tag"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"file_name"`
	LogPath    string `mapstructure:"log_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	DataBase    string `mapstructure:"database"`
	UserName    string `mapstructure:"username"`
	PassWord    string `mapstructure:"password"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"` // 最大空闲连接数
	MaxOpenConn int    `mapstructure:"max_open_conn"` // 最大打开的连接数
	MaxIdleTime int64  `mapstructure:"max_idle_time"` // 连接最大空闲时间

}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	PassWord string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
	DB       int    `mapstructure:"db"`
}

// Init 初始化配置
func Init() (err error) {
	configFile := GetRootDir() + "/config/config.yaml"
	viper.SetConfigFile(configFile) // 指定配置文件（带后缀，可写绝对路径和相对路径两种）
	// 基本上是配合远程配置中心使用的，告诉viper当前的数据使用什么格式去解析
	viper.SetConfigType("yaml") // 远程配置文件传输 确定配置文件的格式
	viper.AddConfigPath(".")    // 指定配置文件的一个寻找路径

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {            // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(globalConfig); err != nil {
		fmt.Println("viper.ReadInConfig() failed")
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}
	// 监控配置文件变化
	viper.WatchConfig()
	// 注意！！！配置文件发生变化后要同步到全局变量Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置信息更新...")
		if err := viper.Unmarshal(globalConfig); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}
	})
	return
}

func GetGlobalConfig() *GlobalConfig {
	return globalConfig
}
