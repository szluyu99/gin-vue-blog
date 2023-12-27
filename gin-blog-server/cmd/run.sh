# 消除 go-sqlite3 的警告
CGO_CFLAGS="-g -O2 -Wno-return-local-addr" go run main.go -m