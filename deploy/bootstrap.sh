#!/bin/bash

# 用 = 而不是 ==: == 是 bash 扩展, 用 sh(dash) 执行时会报错并被当成 false,
# 导致 `sh bootstrap.sh dev` 静默走到生产分支
if [ "$1" = "dev" ]; then
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
# 新版 Docker 的 compose 是插件形式(docker compose), 老版本是独立命令(docker-compose)
if docker compose version > /dev/null 2>&1; then
  docker compose up -d --build
else
  docker-compose up -d --build
fi