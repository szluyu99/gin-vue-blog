#!/usr/bin/env bash
# 本地端到端联调用: 一键启动 / 重启 Redis + 后端 + 博客前台 + 博客后台
# 用法: ./dev.sh [start|stop|restart|status|logs|seed]
# 日志与 pid 放在 .dev/ 下, 已 gitignore

set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RUN_DIR="$ROOT/.dev"
LOG_DIR="$RUN_DIR/log"
REDIS_CONTAINER="gvb-dev-redis"
REDIS_IMAGE="redis:7.0-alpine"

SERVER_PORT=8765
FRONT_PORT=8888
ADMIN_PORT=8889
REDIS_PORT=6379

mkdir -p "$LOG_DIR"

info() { printf '\033[32m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[33m==>\033[0m %s\n' "$*"; }
die() { printf '\033[31m==>\033[0m %s\n' "$*" >&2; exit 1; }

# 端口是否已被监听
port_open() {
  (exec 3<>"/dev/tcp/127.0.0.1/$1") 2>/dev/null && exec 3>&- && return 0
  return 1
}

# 等端口起来, 超时返回 1
wait_port() {
  local port=$1 name=$2 timeout=${3:-60} i=0
  while [ "$i" -lt "$timeout" ]; do
    port_open "$port" && return 0
    i=$((i + 1))
    sleep 1
  done
  warn "$name 等待 ${timeout}s 仍未监听 $port, 看日志: ./dev.sh logs $name"
  return 1
}

pid_file() { echo "$RUN_DIR/$1.pid"; }

alive() {
  local f
  f="$(pid_file "$1")"
  [ -f "$f" ] && kill -0 "$(cat "$f")" 2>/dev/null
}

# setsid 让进程自成进程组, 停止时可以整组杀掉
# (pnpm dev 会派生 vite 子进程, 只杀父进程会留下孤儿占着端口)
# 注意 setsid 自己会 fork, 所以不能直接记录 $!, 那是个转瞬即逝的中间进程;
# 让被拉起的 bash 自己写下 $$ 再 exec, pid 才是真正的服务进程, 也是新进程组的组长
spawn() {
  local name=$1 dir=$2 pf i=0
  shift 2
  pf="$(pid_file "$name")"
  rm -f "$pf"
  # 三处重定向都是必须的, 否则从脚本/工具里调用 ./dev.sh start 会看起来"卡住":
  #   1. 被拉起的服务: stdin 接 /dev/null, stdout/stderr 进日志文件
  #   2. 外层子 shell 自己也要脱开调用方的 stdout —— 实测这个子 shell 会残留下来
  #      (ppid=1, 一直睡着), 手里还攥着调用方的 stdout 管道, 调用方等不到 EOF 就一直挂着
  (cd "$dir" && setsid bash -c 'echo $$ >"$1"; shift; exec "$@"' bash "$pf" "$@" \
    </dev/null >"$LOG_DIR/$name.log" 2>&1 &) </dev/null >/dev/null 2>&1
  while [ ! -s "$pf" ] && [ "$i" -lt 50 ]; do
    i=$((i + 1))
    sleep 0.1
  done
  [ -s "$pf" ] || die "$name 启动失败, 看日志: ./dev.sh logs $name"
}

# 查端口被谁占着, 返回 "pid/命令" 或空串
port_owner() {
  local port=$1 out=''
  if command -v ss >/dev/null 2>&1; then
    out="$(ss -ltnpH "sport = :$port" 2>/dev/null | grep -o 'pid=[0-9]*,fd=[0-9]*' | head -1 | cut -d= -f2 | cut -d, -f1)"
  fi
  if [ -z "$out" ] && command -v lsof >/dev/null 2>&1; then
    out="$(lsof -ti "tcp:$port" -sTCP:LISTEN 2>/dev/null | head -1)"
  fi
  [ -z "$out" ] && return 1
  echo "$out/$(ps -p "$out" -o comm= 2>/dev/null)"
}

# 端口被别的进程占着就直接报错, 否则 vite 会自动换端口, 联调地址对不上
ensure_free() {
  local port=$1 name=$2 owner
  if port_open "$port" && ! alive "$name"; then
    owner="$(port_owner "$port" || true)"
    if [ -n "$owner" ]; then
      die "端口 $port 已被 $owner 占用($name 要用), 处理掉再启动, 或 ./dev.sh stop --force"
    fi
    # 查不到占用者: 端口在本机可见但进程不在当前 PID 命名空间里(例如容器内外混用)
    die "端口 $port 已被占用但查不到进程($name 要用), 可能不在当前 PID 命名空间, 需到对应环境里停掉"
  fi
}

# Redis: 已在跑就复用; 否则优先本机 redis-server, 再退回 docker
start_redis() {
  if port_open "$REDIS_PORT"; then
    info "Redis 已在 $REDIS_PORT 运行, 复用"
    return 0
  fi
  if command -v redis-server >/dev/null 2>&1; then
    info "启动本机 redis-server"
    spawn redis "$ROOT" redis-server --port "$REDIS_PORT"
  elif command -v docker >/dev/null 2>&1; then
    info "本机没有 redis-server, 用 docker 起一个 ($REDIS_IMAGE)"
    docker image inspect "$REDIS_IMAGE" >/dev/null 2>&1 || docker pull "$REDIS_IMAGE" >/dev/null
    docker rm -f "$REDIS_CONTAINER" >/dev/null 2>&1
    docker run -d --name "$REDIS_CONTAINER" -p "$REDIS_PORT:6379" "$REDIS_IMAGE" >/dev/null \
      || die "docker 启动 Redis 失败"
    touch "$RUN_DIR/redis.docker"
  else
    die "需要 Redis 但既没有 redis-server 也没有 docker, 后端会启动失败"
  fi
  wait_port "$REDIS_PORT" redis 30
}

# 后端: 先编译再跑二进制, 这样 pid 就是真正的服务进程, 停止/重启不会留孤儿
# 工作目录必须是 cmd/, 配置里的 ../config.yml、gvb.db、../public/uploaded 都是相对它的
start_server() {
  info "编译后端"
  (cd "$ROOT/gin-blog-server" && go build -o "$RUN_DIR/gvb-server" ./cmd) \
    || die "后端编译失败"
  info "启动后端 :$SERVER_PORT"
  spawn server "$ROOT/gin-blog-server/cmd" "$RUN_DIR/gvb-server"
  wait_port "$SERVER_PORT" server 30
}

# 前端: .env.development 里 VITE_USE_MOCK 默认是 true(纯前端预览用),
# 联调要打到真后端, 所以在命令行覆盖 —— vite 里 shell 环境变量优先级高于 .env 文件
start_web() {
  local name=$1 dir=$2 port=$3
  [ -d "$dir/node_modules" ] || {
    info "$name 安装依赖"
    (cd "$dir" && pnpm install) || die "$name pnpm install 失败"
  }
  info "启动 $name :$port (VITE_USE_MOCK=false)"
  spawn "$name" "$dir" env VITE_USE_MOCK=false pnpm dev
  wait_port "$port" "$name" 60
}

stop_one() {
  local name=$1 f pid
  f="$(pid_file "$name")"
  [ -f "$f" ] || return 0
  pid="$(cat "$f")"
  if kill -0 "$pid" 2>/dev/null; then
    kill -TERM "-$pid" 2>/dev/null || kill -TERM "$pid" 2>/dev/null
    for _ in 1 2 3 4 5; do
      kill -0 "$pid" 2>/dev/null || break
      sleep 1
    done
    kill -KILL "-$pid" 2>/dev/null
    info "已停止 $name (pid $pid)"
  fi
  rm -f "$f"
}

do_stop() {
  stop_one admin
  stop_one front
  stop_one server
  stop_one redis
  # pid 文件丢了(手动删过, 或进程是别的环境起的)时按端口兜底,
  # 会杀掉占着这四个端口的任意进程, 所以要显式加 --force
  if [ "${1:-}" = "--force" ]; then
    local owner pid
    for port in "$ADMIN_PORT" "$FRONT_PORT" "$SERVER_PORT"; do
      port_open "$port" || continue
      owner="$(port_owner "$port" || true)"
      pid="${owner%%/*}"
      if [ -n "$pid" ]; then
        kill -TERM "-$pid" 2>/dev/null || kill -TERM "$pid" 2>/dev/null
        sleep 1
        kill -KILL "-$pid" 2>/dev/null
        info "已强制停止占用 $port 的进程 $owner"
      else
        warn "端口 $port 仍被占用但查不到进程, 可能不在当前 PID 命名空间"
      fi
    done
  fi
  if [ -f "$RUN_DIR/redis.docker" ]; then
    docker rm -f "$REDIS_CONTAINER" >/dev/null 2>&1 && info "已停止 Redis 容器"
    rm -f "$RUN_DIR/redis.docker"
  fi
}

# 初始化基础数据(菜单/资源/角色/默认用户/网站配置), 可重复执行, 已存在的会跳过
do_seed() {
  port_open "$REDIS_PORT" || start_redis
  info "初始化基础数据"
  (cd "$ROOT/gin-blog-server/cmd" && sh generate_data.sh) || die "初始化基础数据失败"
  info "默认账号: admin / 123456, guest / 123456"
}

do_start() {
  command -v go >/dev/null 2>&1 || die "未找到 go"
  command -v pnpm >/dev/null 2>&1 || die "未找到 pnpm"

  local first_run=0
  [ -f "$ROOT/gin-blog-server/cmd/gvb.db" ] || first_run=1

  start_redis
  if [ "$first_run" = 1 ]; then
    warn "没有 gvb.db, 视为首次启动, 先建库再灌基础数据"
    start_server
    stop_one server
    do_seed
  fi
  ensure_free "$SERVER_PORT" server
  ensure_free "$FRONT_PORT" front
  ensure_free "$ADMIN_PORT" admin
  start_server
  start_web front "$ROOT/gin-blog-front" "$FRONT_PORT"
  start_web admin "$ROOT/gin-blog-admin" "$ADMIN_PORT"

  echo
  info "全部就绪"
  echo "  博客前台  http://localhost:$FRONT_PORT"
  echo "  博客后台  http://localhost:$ADMIN_PORT   (admin / 123456)"
  echo "  后端接口  http://localhost:$SERVER_PORT/api"
  echo "  Swagger   http://localhost:$SERVER_PORT/swagger/index.html"
  echo "  日志      ./dev.sh logs [server|front|admin|redis]"
}

do_status() {
  local name port
  for entry in "redis:$REDIS_PORT" "server:$SERVER_PORT" "front:$FRONT_PORT" "admin:$ADMIN_PORT"; do
    name="${entry%%:*}"
    port="${entry##*:}"
    if port_open "$port"; then
      if alive "$name"; then
        printf '  %-7s 运行中  端口 %-5s pid %s\n' "$name" "$port" "$(cat "$(pid_file "$name")")"
      elif [ "$name" = redis ] && [ -f "$RUN_DIR/redis.docker" ]; then
        printf '  %-7s 运行中  端口 %-5s docker 容器 %s\n' "$name" "$port" "$REDIS_CONTAINER"
      else
        printf '  %-7s 运行中  端口 %-5s (非本脚本启动)\n' "$name" "$port"
      fi
    else
      printf '  %-7s 未运行  端口 %s\n' "$name" "$port"
    fi
  done
}

case "${1:-start}" in
  start) do_start ;;
  stop) do_stop "${2:-}" ;;
  restart)
    do_stop "${2:-}"
    do_start
    ;;
  status) do_status ;;
  seed) do_seed ;;
  logs)
    name="${2:-server}"
    [ -f "$LOG_DIR/$name.log" ] || die "没有 $name 的日志"
    tail -f "$LOG_DIR/$name.log"
    ;;
  *)
    echo "用法: ./dev.sh [start|stop|restart|status|seed|logs <name>]"
    echo "      stop / restart 可加 --force: pid 文件丢了时按端口强杀占用进程"
    exit 1
    ;;
esac
