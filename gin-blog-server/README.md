# gin-blog-server

博客后端。Go 1.21 + Gin + GORM，支持 SQLite / MySQL，**Redis 为必需依赖**（默认 `127.0.0.1:6379` 的 DB 7）。

监听 `:8765`，接口前缀 `/api`。

## 启动

配置文件是 `config.yml`（程序从 `cmd` 目录下以 `../config.yml` 读取）。默认为 SQLite，开箱即用：

```bash
go mod tidy
cd cmd
go run main.go
```

`DbAutoMigrate: true` 时启动会自动迁移表结构，但**只建表、不造数据**。

## 初始化数据

### SQLite

第一次启动后需要手动生成基础数据（菜单、资源、角色、默认用户、网站配置、页面封面）：

```bash
cd cmd
sh generate_data.sh          # 内部执行 go run main.go -t "all"
```

生成默认角色 `admin`、`guest` 和同名用户，密码都是 `123456`。后台登录用 `admin / 123456`。

如需一个拥有全部权限的超级管理员（`superadmin / superadmin`）：

```bash
sh create_superadmin.sh
```

> `generate_data.sh` 只生成系统基础数据，**不含文章、分类、标签、评论、友链等内容数据**，这些需要自己在后台添加。
> `assets/gvb.sql` 里有一份示例内容数据，但它是 MySQL 格式（Navicat 导出），SQLite 下需要先转换才能导入。

### MySQL

`config.yml` 中改 `DbType: "mysql"` 并填好连接信息，然后导入 `assets/gvb.sql`（已包含表结构和示例数据）。

## 需要注意的缓存

页面封面等数据会缓存到 Redis 且**没有过期时间**。如果直接改了数据库但接口仍返回旧值，删掉对应缓存即可：

```bash
redis-cli -n 7 del page
```

## 前端不想启动后端？

博客前台和后台都内置了 Mock 模式，把对应目录 `.env.development` 中的 `VITE_USE_MOCK` 设为 `true` 即可完全脱离后端运行，详见各自的 README。

## 其他

```bash
cd cmd
sh run_swag.sh    # swag init 生成 Swagger 文档, 之后会直接启动服务
```

目录约定：`internal/model` 返回 error，`internal/handle` 负责错误码与响应；JSON 字段统一 **小写 + 下划线**，Go 结构体用驼峰。
