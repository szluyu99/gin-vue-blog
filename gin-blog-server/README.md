## 重构

- 使用 slog 作为日志库。对于 Go 1.21 以后版本的项目，slog 内置在标准库中。对于之前版本，需要 引入 `exp/slog` 来使用（TODO: 日志持久化）
- 修改为 Golang 官方标准项目结构 
- 使用 Swagger 生成接口文档（TODO: 给每个接口完善注释，而且感觉 go swagger 并不好用，寻找更好的方案）
- 修改全局变量的使用为依赖注入形式, 统一在 main.go 初始化
- 修改 utils 包架构, 精简项目架构
- 开发环境使用 sqlite 进行单元测试, 线上支持 MySQL 环境
- 移除 Casbin，自己实现基于 RBAC 的权限控制
- 之前版本错误码使用 panic 机制, 使用 gin 中间件全局捕获 panic
- TODO: 使用 errors 机制来统一错误码机制

新版本项目运行：运行/初始化数据 相关文件都在 cmd 下

1、MySQL，`config.yml` 中 `DbType = "mysql"` 并修改 MySQL 连接信息， 导入 `assets/gvb.sql`（也可以手动初始化数据，参考 SQLite 的情况）

2、SQLite，`config.yml` 中 `DbType = "sqlite"`，并执行初始化数据操作：
- `create_superadmin.sh` 创建一个超级管理员（拥有所有权限），用户名 superadmin 密码 superadmin
- `generate_data.sh` 初始化 三个默认角色 (admin, user, guest) + 三个默认用户 (admin, user, guest, 密码都是 123456),  初始化配置数据, 初始化页面数据, 初始化资源数据 (TODO: )

---

以下是一些旧版本的笔记记录，可以不看

## 项目部署

项目使用 docker-compose 一键部署运行。

./config/config.docker.toml 是部署时读取的配置文件，其中有些变量会由 docker-compose 中的环境变量覆盖

具体可参考 deploy/start/docker-compose.yml

配置文件读取值的优先级: 环境变量 > config.docker.toml 中的值，因此使用 docker-compose 运行时，修改其中的 environment 即可。

## 开发规范

Model 层返回 error, Service 层返回 code, Controller 层通过 `GetMsg(code)` 获取具体信息返回给前端

- 错误码在 global/errmsg 中维护

JSON 传输数据一律使用 **小写 + 下划线**，Go 项目中使用 **驼峰**

## 数据库

MySQL 设置字段 true false 使用 tinyint, 注意在 Golang 定义结构体时该字段要使用指针
否则结构体的字段会被初始化其零值 (0), 设置为指针可以用 nil 来判断

## 单元测试

单元测试很重要，尤其是项目开发到中后期要进行重构和调整，需要使用测试来减少每次手动调用接口花费的时间，同时可以保证程序的正确性。

## Gin

### gin 的参数校验 validator 中零值问题

gin 的参数校验是基于 validator 的，如果给了 `required` 标签，则不能传入零值
- 比如字符串的不能传入空串
- int 类型的不能传入 0
- bool 类型的不能传入 false。

有时候需要参数是必填，并且需要可以传入 0，例如 sex 为 0 表示女，为 1 表示男，且必填
通过定义 int 类型的**指针** 解决这个问题：指针的零值是 `nil`

```go
Golang 解析 JSON 格式数据时，若以 any 接收数据，则会按照以下规则解析:
bool, for JSON booleans
float64, for JSON numbers
string, for JSON strings
[]interface{}, for JSON arrays
map[string]interface{}, for JSON objects
nil for JSON null
```

### 关于使用 POST 还是 PUT？

POST:
- 用于提交请求，可以更新或者创建资源，是非幂等的
- 在用户注册时，每次提交都是创建一个用户账号，此时用 POST

PUT:
- 用于向指定的 url 传送更新资源，是幂等的
- 比如修改密码，虽然提交的还是账户名和密码，但是每次提交都只是更新该用户密码，每次请求都只是覆盖原型的值，此时用 PUT

### 菜单、资源是树形的数据，如何处理？

三种方式处理：
1. MySQL 查询时直接查出树形结构（使用自定义函数 或 其他方式）
2. 程序中使用递归构建树形结构数据
3. 通过一个 map，遍历一次就可以完成树的构建 :star:
> map 的方法似乎只能构建出二级菜单？？？ TODO: 待查询资料

### 日志记录

启动相关的日志直接输出到控制台，运行过程中的日志记录到文件


# Nginx 部署

Nginx 配置 Https 参考: [Nginx 服务器 SSL 证书安装部署](https://cloud.tencent.com/document/product/400/35244)
