## Mini_tiktok
项目来源于第六届字节跳动后端训练营，参考往届实现，对照2023.8最新的[极简版抖音APP接口文档](https://apifox.com/apidoc/shared-09d88f32-0b6c-4157-9d07-a36d32d7a75c/api-50707523)进行了完善。

## 技术选型

- [x] Go 1.20
- [x] HTTP框架：Gin
- [x] ORM框架：Gorm
- [x] RPC框架：gRPC
- [x] 数据库：MySQL
- [x] 缓存：Redis
- [x] 服务注册与发现：Consul
- [x] 日志：Zap
- [x] 配置文件：Viper
- [x] 存储：Minio

## 目录结构

```
├── cmd 项目启动入口
│   ├── gatewaysvr 网关服务
│   ├── usersvr 用户服务
│   ├── videosvr 视频服务
│   ├── favoritesvr 点赞服务
│   ├── relationsvr 关注服务
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
├── pkg proto文件
├── script 快速启动脚本
    └── build_all.py 编译所有微服务
    └── server_all.py 启动/停止所有微服务
    └── init_db.sql 数据库初始化
    └── redis.sh 初始化
```

## 项目总结与反思

- [x] 上传视频时间过长，影响用户体验，需要优化
- [x] 后续响应消息体应该具体化，不应该只返回0和1

