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

# 生成本机独有的签名密钥 (JWT / session), 只在文件不存在时生成
# config.docker.yml 里是 release 模式, 后端启动时会拒绝空值和仓库里的示例密钥
SECRET_FILE=start/.env.secrets
if [ ! -f "$SECRET_FILE" ]; then
  gen_secret() {
    if command -v openssl > /dev/null 2>&1; then
      openssl rand -hex 32
    else
      head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n'
    fi
  }
  {
    echo "# 本机自动生成的签名密钥, 不要提交到仓库, 删掉后 bootstrap.sh 会重新生成"
    echo "# 注意: 更换后所有已签发的 token 和 session 立即失效"
    echo "JWT_SECRET=$(gen_secret)"
    echo "SESSION_SALT=$(gen_secret)"
  } > "$SECRET_FILE"
  echo "已生成 $SECRET_FILE"
fi

# 启动新容器
cd start
# 新版 Docker 的 compose 是插件形式(docker compose), 老版本是独立命令(docker-compose)
if docker compose version > /dev/null 2>&1; then
  docker compose up -d --build
else
  docker-compose up -d --build
fi