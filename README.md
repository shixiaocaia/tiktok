## Mini_tiktok
项目来源于第六届字节跳动后端训练营，参考往届实现，对照2023.8最新的[极简版抖音APP接口文档](https://apifox.com/apidoc/shared-09d88f32-0b6c-4157-9d07-a36d32d7a75c/api-50707523)进行了完善。

## 技术选型

- [x] 语言：Go 1.19以上
- [x] HTTP框架：Gin
- [x] ORM: Gorm
- [x] 服务注册与发现：Consul
- [x] 服务间调用：gRPC
- [x] 数据库：MySQL
- [x] 缓存：Redis
- [x] 分布式锁：RedSync
- [x] 配置：Viper
- [x] 日志：Zap
- [x] JWT：jwt-go
- [x] 代码生成：protoc-gen-go

## 目录结构

```
├── cmd 项目启动入口
│   ├── gatewaysvr 网关服务
│   ├── usersvr 用户服务
│   ├── videosvr 视频服务
│   ├── favoritesvr 点赞服务
│   ├── relationsvr 评论服务
│   ├── commentsvr 评论服务
│   ├── messagesvr 聊天服务
│   ├── 具体文件
│        └── config 配置文件
│        └── constant 常量值
│        └── dao 数据库操作
│        └── log 日志配置
│        └── middleware 中间件
│        └── response 响应
│        └── service 服务具体逻辑
│        └── utils 工具类
├── log 日志文件
├── pkg proto文件
├── script 快速启动脚本
├── README.md
```

## 快速开始

```shell
cd script
# 权限问题
sudo chmod -R 777 /home/gopath/src/tiktok/cmd
sudo chmod -R 777 /home/gopath/src/tiktok/script
# 编译所有微服务
python build_all.py "" / gatewaysvr...
# 启动/停止所有微服务
python server_all.py start/stop
```
