## Mini_tiktok

## 技术选型

- [x] 语言：Go 1.19以上
- [x] HTTP框架：Gin
- [x] ORM: Gorm
- [x] 数据库：MySQL
- [x] 缓存：Redis
- [x] 配置：Viper
- [x] 日志：Zap
- [x] JWT：jwt-go
- [x] 代码生成：protoc-gen-go

## 目录结构

```
├── cmd 项目启动入口
│   ├── gatewaysvr 网关服务
│   ├── testsvr 测试文件
│   ├── usersvr 用户服务
│        └── config 配置文件
│        └── constant 常量值
│        └── dao 数据库操作
│        └── log 日志配置
│        └── middleware 中间件
│        └── response 响应
│        └── service 服务具体逻辑
│        └── utils 工具类
│        └── main.go 项目启动入口
├── log 日志文件
├── pkg proto文件及代码生成
├── model 数据库模型
├── README.md