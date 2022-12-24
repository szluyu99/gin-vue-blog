# 二次开发指南

后端服务的 Dockerfile 参考 gin-blog-server 目录

直接在修改 gin-blog-server 中的后端源码

然后在 start 目录重新用 docker-compose 构建运行

```bash
docker-compose up -d --build
```
