FROM golang:1.21-alpine as BUILDER
WORKDIR /gvb

COPY go.mod go.mod
COPY go.sum go.sum
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download
COPY . .
RUN cd cmd && go build -o server . 

FROM alpine:3.19
ENV WORK_PATH /gvb
WORKDIR ${WORK_PATH}
COPY --from=0 ${WORK_PATH}/cmd/server .
COPY --from=0 ${WORK_PATH}/config.docker.yml .
COPY --from=0 ${WORK_PATH}/assets/ip2region.xdb ./assets/ip2region.xdb

EXPOSE 8765

ENTRYPOINT ./server -c config.docker.yml