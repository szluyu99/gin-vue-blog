# 进入根目录, 生成 swagger 文档
cd ..
swag init -g ./cmd/main.go

echo "------swag init done------"

# 回到 cmd 目录, 运行 main.go
cd cmd
# 消除 go-sqlite3 的警告
CGO_CFLAGS="-g -O2 -Wno-return-local-addr" go run main.go