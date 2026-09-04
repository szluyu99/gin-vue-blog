cd generate-data

# all | config | auth | page | resource | demo
# all 生成所有信息
# config 生成配置信息
# auth 生成默认角色 admin, guest, 以及对应的默认用户 admin, guest
# page 生成默认页面信息
# resource 生成默认资源信息
# demo 生成本地测试用的样例内容(分类/标签/文章/评论/留言/友链), 库里已有文章则跳过
go run main.go -t "all"