最新版：有以下几个脚本
- `build_web.sh`：本地打包 Web 项目（需要有 node 环境）并将静态资源移到容器构建目录
- `clean_docker.sh`：清理本项目相关的旧 Docker 容器
- `bootstrap.sh`：使用 docker node 打包静态资源
- `bootstrap.sh dev`：使用本机的 pnpm 打包前端静态资源

一般来说，直接运行 `bootstrap.sh` 即可，每次自动清理旧容器，打包最新代码，并构建新容器。第一次会比较耗时，但是后面会有缓存就会快很多

---

一键运行分成两步：
1. 环境准备: 需要 Docker 和 Docker Compose 环境
2. 开始运行: 执行 `bootstrap.sh` 脚本
3. 运行可能遇到的问题: 如果运行失败了，可以来这看看

# 1、环境准备

## Windows 环境 / Mac 环境

> 相信大部分小伙伴日常都是使用 Windows, 可以在 Windows 上使用 Docker 一键运行项目看效果。
> 
> 但是如果真的要部署到云服务器上，还是建议使用 Linux 哈。

安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/) 即可，安装完自带 Docker 和 Docker Compose。

> 软件安装无脑下一步就行，遇到问题多百度。

## Linux 环境

> 如果是桌面版 Linux, 也可以直接安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/) 即可

> 本文档编写的标准是基于 **Ubuntu 系统**，如使用其他 Linux 系统可能会有些许差异，遇到问题多百度。

#### 1. 安装必备依赖

```bash
sudo apt update && sudo apt install -y vim curl git
```

#### 2. 安装 Docker

> 以下内容摘自官方文档，且最新安装方法永远参考官方文档：[Install Docker Engine on Ubuntu](https://docs.docker.com/engine/install/ubuntu/)

1. 使用 `apt` 更新并安装依赖包

```bash
sudo apt-get update

sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
```

2. 添加 Docker 官方 GPG 密钥

```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
```

3. 设置仓库

```bash
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

4. 安装 Docker Engine

```bash
sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin
```

5. 验证是否安装成功

```
sudo docker run hello-world
```

#### 3. 安装 Docker Compose

> 最新参考官方文档：[Install the Compose standalone](https://docs.docker.com/compose/install/other/)

1. 下载

```bash
sudo curl -SL https://github.com/docker/compose/releases/download/v2.14.2/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
# 上面的速度太慢可以尝试换源
# sudo curl -SL https://get.daocloud.io/docker/compose/releases/download/v2.14.2/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
```

2. 授权

```bash
sudo chmod +x /usr/local/bin/docker-compose
```

# 2、开始运行

## 不修改源码，直接快速运行

```bash
# 拉取项目
git clone https://github.com/szluyu99/gin-vue-blog
cd gin-vue-blog

cd deploy
./bootstrap.sh
```

如果只是查看运行效果，可以直接按照以上操作进行，采用默认环境变量即可。

如果是线上部署，建议修改 start 中 .env 配置文件中的重要参数（数据库密码等）

## 修改源码后，进行运行

修改 gin-blog-server 后端源码，不需要额外做什么，因为后端构建镜像时直接依赖 gin-blog-server 中的 Dockerfile

修改 gin-blog-admin 后台源码后，需要将打包后静态资源覆盖 build/web/html/admin 

修改 gin-blog-front 前台源码后，需要将打包后静态资源覆盖 build/web/html/blog

执行以上操作后，在 start 目录重新构建运行：

```bash
docker-compose up -d --build
```

以上操作我放在 `build_web.sh` 中了，直接执行该脚本即可

## 线上部署的注意事项

线上部署，其实就是在服务器上运行起这个项目，有两个地方的参数推荐调整一下。

> 如果你不在乎安全性，可以忽略这些建议，直接部署项目。

1、执行 `docker-compose up -d` 前建议修改 `.env` 文件中的环境变量

为了保证安全性，以下几项建议修改：
- `REDIS_PASSWORD` Redis 连接密码 
- `MYSQL_ROOT_PASSWORD` MySQL 连接密码

其他根据需求修改，一般不用动

2、后端镜像的构建直接依赖于 gin-blog-server 中的源码，构建时加载的是 config/config.docker.toml 配置文件。

其中有一些参数推荐重新设置，具体可以查看该文件注释。

---

注意，默认管理员用户为 admin 123456。<font color=red> 项目运行成功后，务必使用管理员登录后台，修改其登录密码。</font>

# 3、运行可能遇到的问题

## gvb-mysql 和 gvb-server 无法启动

gvb-mysql 服务 和 gvb-server 服务运行不起来，很有可能是因为你本机已经有了一个 MySQL 服务占用了 3306 端口，所以 gvb-mysql 运行失败，然后由于 gvb-server 依赖 gvb-mysql，所以它也运行失败。
- 解决方案一：关闭本机的 MySQL 服务，自行百度
- 解决方案二：修改 .env 文件中的 `MYSQL_PORT` 为其他端口

> 解决方案二中修改了 MySQL 的端口是 Docker 容器对外暴露的端口，即我如果修改为 33069，则我本机通过 Navicat 连接的应该是 127.0.0.1:33069 的 mysql 服务。但是 MySQL 容器内是固定运行在 3306 端口的。但是容器和主机，以及容器之前是互相隔离的，所以不会有影响。


## gvb-web 和 gvb-server 无法启动

gvb-web 的日志中报错如下：

```
'：No such file or directory
```

如果是 Windows 系统，需要先执行以下指令，否则 Docker 过程会出 BUG。

或者直接下载 ZIP 而不是通过 git clone 克隆项目。

Linux 和 Mac 不需要进行该操作。

> 原因是该项目开发时基于 Linux，本项目规范使用 lf 换行符。而 Windows 的 git 在自动拉取项目时会将项目文件中换行符转换为 ctlf，经过测试，构建过程会产生 BUG。

```bash
# 防止 git 自动将换行符转换为 crlf
git config --global core.autocrlf false
```

## 如果已经运行过，又修改了 .env 中的数据库密码重新运行

需要把原来的数据文件删除，即 `start/gvb` 目录