#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import sys

# 设置允许核心转储文件
os.system("ulimit -c unlimited")

# 定义合法的操作和服务器名称列表
valid_actions = ["start", "stop"]
valid_server_names = ["gatewaysvr", "commentsvr", "favoritesvr", "relationsvr", "usersvr", "videosvr", "messagesvr", "all"]

# 获取命令行参数
if len(sys.argv) < 3:
    print("Usage: ./server.py [start|stop] [gatewaysvr|commentsvr|favoritesvr|relationsvr|usersvr|videosvr|all|...]")
    sys.exit(1)

action = sys.argv[1]
server_name = sys.argv[2]

print("ACTION=" + action)
print("SERVER_NAME=" + server_name)

# 定义一个函数来执行指定的操作和服务器名称
def run(action, server_name):
    if action not in valid_actions:
        print("param1 should be start or stop")
        sys.exit(1)
    else:
        if server_name not in valid_server_names:
            print("server name is not correct")
            sys.exit(1)
        elif server_name == "all":
            # 执行所有服务器的操作
            server_names = valid_server_names[:-1]  # 除了 "all"
            script_path = os.getcwd()
            print("SCRIPT_PATH====" + script_path)
            for srv_name in server_names:
                script_folder = os.path.join(script_path, "../cmd/" + srv_name + "/script")
                os.chdir(script_folder)
                os.system("./server.sh " + action)
        else:
            # 执行单个服务器的操作
            script_folder = os.path.join(os.getcwd(), "../cmd/" + server_name + "/script")
            os.chdir(script_folder)
            os.system("./server.sh " + action)

run(action, server_name)
