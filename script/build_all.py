# -*- coding: utf-8 -*-

import os
import sys

# 定义要进入的文件夹列表
folders = ["gatewaysvr", "commentsvr", "favoritesvr", "relationsvr", "usersvr", "videosvr", "messagesvr"]

# 获取命令行参数
if len(sys.argv) > 1:
    selected_command = sys.argv[1]
else:
    selected_command = None

# 遍历文件夹列表
for folder in folders:
    # 构建文件夹路径
    folder_path = os.path.join("/home/ubuntu/go/src/tiktok/cmd", folder, "script")

    # 切换到文件夹路径
    os.chdir(folder_path)

    # 执行build.sh文件
    if selected_command == folder or selected_command == None:
        os.system("sh build.sh")

