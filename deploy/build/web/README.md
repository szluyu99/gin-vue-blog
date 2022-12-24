# 二次开发指南

修改 gin-blog-admin 和 gin-blog-front 源码后，用以下指令打包会生成 dist 目录

```bash
pnpm build
```

将 gin-blog-admin 打包生成的 dist 中文件, 移到当前目录 html/admin 下

将 gin-blog-front 打包生成的 dist 中文件, 移到当前目录 html/blog 下

然后在 start 目录重新用 docker-compose 构建运行

```bash
docker-compose up -d --build
```
