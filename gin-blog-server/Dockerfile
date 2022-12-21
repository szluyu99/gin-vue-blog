FROM golang:alpine as builder
# 设置工作目录
WORKDIR /app
# 将当前目录内容拷到工作目录 (相对路径)
COPY . .
# 配置 golang 环境
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest
WORKDIR /app
COPY --from=0 /app/server .
COPY --from=0 /app/config/config.docker.toml .
COPY --from=0 /app/config/ip2region.xdb ./config/ip2region.xdb

# 后台接口
EXPOSE 8765 
# 前台接口
EXPOSE 5678
ENTRYPOINT ./server -c config.docker.toml