# 二次开发指南

后端服务的 Dockerfile 参考 gin-blog-server 目录

直接在修改 gin-blog-server 中的后端源码，然后执行 `./bootstrap.sh`

容器的启动脚本是 `gin-blog-server/docker-entrypoint.sh`：先执行 `generate-data`
初始化角色/接口资源/菜单（幂等，以代码中的 `internal/model/seed_resource.go` 为准），
再启动 `server`。所以新增后台接口后，改完 `seed_resource.go` 重启容器即可生效。