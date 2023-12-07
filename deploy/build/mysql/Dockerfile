FROM mysql:8.0

# 定义工作目录
ENV WORK_PATH /usr/local/work
# 定义被容器自动执行的目录
ENV AUTO_RUN_DIR /docker-entrypoint-initdb.d
# 定义要执行的 shell 文件
ENV RUN_SHELL run.sh

COPY ./mysql.cnf /etc/mysql/conf.d/
# 把数据库初始化数据的文件复制到工作目录
COPY ./gvb.sql ${WORK_PATH}/
# 把执行初始化的脚本放到 /docker-entrypoint-initdb.d 目录下
COPY ./${RUN_SHELL} ${AUTO_RUN_DIR}/

# 给执行文件添加可执行权限
RUN chmod a+x ${AUTO_RUN_DIR}/${RUN_SHELL}