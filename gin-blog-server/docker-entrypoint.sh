#!/bin/sh
# 启动前初始化种子数据: 资源/菜单/角色 以代码 (internal/model/seed_resource.go) 为准,
# 重复执行是幂等的 (靠唯一约束跳过), 因此每次启动都跑一次, 新增接口重启即生效
set -e

./generate-data -c config.docker.yml -sqlite-parent=false

exec ./server -c config.docker.yml
