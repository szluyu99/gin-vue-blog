FROM golang:1.20-alpine as BUILDER
WORKDIR /gvb

COPY go.mod go.mod
COPY go.sum go.sum
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download
COPY . .
RUN go build -o server . 

FROM alpine:3.19
ENV WORK_PATH /gvb
WORKDIR ${WORK_PATH}
COPY --from=0 ${WORK_PATH}/server .
COPY --from=0 ${WORK_PATH}/config/config.docker.toml .
COPY --from=0 ${WORK_PATH}/assets/ip2region.xdb ./assets/ip2region.xdb

# COPY --from=0 ${WORK_PATH}/assets/wait-for .
# RUN chmod a+x ./wait-for

# 后台接口
EXPOSE 8765
# 前台接口
EXPOSE 5678

ENTRYPOINT ./server -c config.docker.toml

###############################################
# CMD sleep 5 && ./server -c config.docker.toml