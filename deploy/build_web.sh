#!/bin/bash

# 未安装 pnpm 则安装
if ! command -v pnpm &> /dev/null; then
  npm install -g pnpm
else 
  echo "pnpm 已安装"
fi

# 打包 gin-blog-front 并移到 docker 部署目录
cd ../gin-blog-front
pnpm install
pnpm build
rm -rf ../deploy/build/web/dist_blog
mv ./dist ../deploy/build/web/dist_blog

# 打包 gin-blog-admin 并移到 docker 部署目录
cd ../gin-blog-admin
pnpm install
pnpm build
rm -rf ../deploy/build/web/dist_admin
mv ./dist ../deploy/build/web/dist_admin