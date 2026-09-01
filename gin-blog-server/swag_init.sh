#!/bin/sh
# 生成 Swagger 文档到 docs/, 生成结果需要一起提交 (CI 会校验是否最新)
# 没有 swag 命令时先安装: go install github.com/swaggo/swag/cmd/swag@v1.16.6
swag init -g ./cmd/main.go
