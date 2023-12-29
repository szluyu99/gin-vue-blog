cd generate-data

# all | config | auth | page | resource
# config 生成配置信息
# auth 生成 3 个默认角色 admin, user, guest, 以及 3 个默认用户 admin, user, guest
# page 生成默认页面信息
# resource 生成默认资源信息
go run main.go -type "all"