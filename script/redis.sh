#!/bin/bash

# 入 口 函 数
function Init() {


        GetVar
        WriteToConf

        docker run -d --network host --restart unless-stopped --name $name -v /docker/redis/$name/conf:/etc/redis/  -v $name-data:/data cjm/redis redis-server /etc/redis/redis.conf
        # 软 链 接
        ln -s /var/lib/docker/volumes/$name-data/_data /docker/redis/$name/data
}

# 获 取 变 量
function GetVar() {
        echo "正 在 使 用 redis master节 点 脚 本 ！ "
        read -t 20 -p "请 输 入 容 器 映 射 到 主 机 的 端 口 （ 容 器 内 默 认 为 6379， 主 机 默 认 为 6379） ： " port
        port=${port:-6379}
        echo "主 机 端 口 为 ${port}"

        read -t 20 -p "请 输 入 容 器 运 行 时 的 容 器 名 字 （ 默 认 为 redis） ： " name
        name=${name:-redis}
        echo "容 器 名 字 为 ${name}"

        read -t 20 -p "请 输 入 主 节 点 密 码 （ 默 认 为 123456） ： " password
        password=${password:-123456}
        echo "主 节 点 密 码 为 ${password}"

        # 创 建 需 要 映 射 的 目 录
        mkdir -p /docker/redis/$name/conf /docker/redis/$name/data
}


# 编 写 配 置 文 件
function WriteToConf() {
cat > /docker/redis/$name/conf/redis.conf << EOF

# 容 器 内 部 端 口
port $port
# daemonize yes

logfile "redis-$port.log"
dbfilename "dump-redis-$port.rdb"

# 访 问 redis-server密 码
requirepass $password

# 主 节 点 密 码 （ 也 需 要 配 置 ， 不 然 当 master无 法 成 为 新 的 master节 点 的 slave节 点 ）
masterauth $password

# 设 置 redis最 大 内 存
maxmemory 500MB
# 设 置 redis内 存 淘 汰 策 略
maxmemory-policy volatile-lru

# 开 启 AOF
appendonly yes
# 总 是 追 加 到 AOF文 件
appendfsync always
# 超 过 100MB就 进 行 重 写
auto-aof-rewrite-min-size 100mb
# 超 过 增 长 率 就 进 行 重 写
auto-aof-rewrite-percentage 40
# 开 启 AOF-RDB 混 合 持 久 化
aof-use-rdb-preamble yes

# 配 置 慢 查 询 日 志 （ 不 超 过 100微 秒 就 不 会 被 记 录 ）
slowlog-log-slower-than 100
# 最 多 记 录 100条 ， 先 进 先 出
slowlog-max-len 100
EOF
}

Init