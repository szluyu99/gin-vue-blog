FROM golang:1.26-alpine as BUILDER
WORKDIR /gvb

COPY go.mod go.mod
COPY go.sum go.sum
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download
COPY . .
RUN cd cmd && go build -o server .
# 权限数据(资源/菜单/角色)以代码为准, 由 generate-data 初始化, 见 docker-entrypoint.sh
RUN cd cmd/generate-data && go build -o /gvb/cmd/generate-data-bin .

FROM alpine:3.19
ENV WORK_PATH /gvb
WORKDIR ${WORK_PATH}
COPY --from=0 ${WORK_PATH}/cmd/server .
COPY --from=0 ${WORK_PATH}/cmd/generate-data-bin ./generate-data
COPY --from=0 ${WORK_PATH}/config.docker.yml .
COPY --from=0 ${WORK_PATH}/assets/ip2region.xdb ./assets/ip2region.xdb
COPY --from=0 ${WORK_PATH}/assets/templates ./assets/templates
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod a+x /docker-entrypoint.sh
EXPOSE 8765

ENTRYPOINT ["/docker-entrypoint.sh"]