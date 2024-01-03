cd generate-data

# all | config | auth | page | resource
# all 生成所有信息
# config 生成配置信息
# auth 生成默认角色 admin, guest, 以及对应的默认用户 admin, guest
# page 生成默认页面信息
# resource 生成默认资源信息
go run main.go -t "all"