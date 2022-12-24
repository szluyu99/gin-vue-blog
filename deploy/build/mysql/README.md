# 二次开发指南

mysql 镜像的主要作用是初始化数据库中数据，核心在于 gvb.sql 文件

这是启动 mysql 容器后会自动执行的 sql 文件

如需更改数据库初始数据，修改 gvb.sql 文件

然后在 start 目录重新用 docker-compose 构建运行

```bash
docker-compose up -d --build
```
