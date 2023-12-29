# 进入根目录, 生成 swagger 文档
cd ..
swag init -g ./cmd/main.go

echo "------swag init done------"

# 回到 cmd 目录, 运行 main.go
cd cmd
go run main.go