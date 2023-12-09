#!/bin/bash
# 清理本项目构建的 Docker 容器

stop_and_remove_container() {
  container_name=$1
  if docker ps -a --format '{{.Names}}' | grep -q "$container_name"; then
    docker stop $container_name
    docker rm $container_name
    echo "已经停止并删除容器 $container_name"
  else
    echo "未找到容器 $container_name"
  fi
}

stop_and_remove_container "gvb-web"
stop_and_remove_container "gvb-server"
stop_and_remove_container "gvb-redis"
stop_and_remove_container "gvb-mysql"

echo "本项目相关 Docker 容器清理完成"