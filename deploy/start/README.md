建议执行 deploy 目录下的 `bootstrap.sh` 脚本，会做一些清理旧容器等功能

也可以进入 start 目录后，执行以下命令：

```bash
docker compose up -d --build
# 旧版 Docker 没有 compose 插件时用: docker-compose up -d --build
```
