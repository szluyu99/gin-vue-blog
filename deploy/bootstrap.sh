#!/bin/bash

if [ "$1" == "dev" ]; then
  # 开发环境, 使用本机的 pnpm 环境打包 (相对于 start/Dockerfile 的路径)
  export WEB_BUILD_CONTEXT="../build/web"
  ./build_web.sh
else
  # 生产环境, 使用 docker 容器的 node 打包
  export WEB_BUILD_CONTEXT="../.."
fi

# 清理旧容器
./clean_docker.sh

# 启动新容器
cd start
docker-compose up -d --build