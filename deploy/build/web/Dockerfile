FROM nginx:alpine

# 将 Nginx 配置文件模板拷到容器中
COPY default.conf.template /etc/nginx/conf.d/default.conf.template
COPY default.conf.ssl.template /etc/nginx/conf.d/default.conf.ssl.template

# 博客前台静态资源 拷到容器
ADD html/blog/ /usr/share/nginx/html/
# 博客后台静态资源 拷到容器
ADD html/admin/ /usr/share/nginx/html/admin

# 初始化脚本, 根据环境变量和模板生成 Nginx 配置文件
COPY ./run.sh /docker-entrypoint.sh
RUN chmod a+x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]

# 每次容器启动时执行
CMD [ "nginx", "-g", "daemon off;" ]

EXPOSE 80
EXPOSE 443