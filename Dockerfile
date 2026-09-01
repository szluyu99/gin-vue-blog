# 阶段一: 打包前后台静态资源
# Vite 8 要求 Node >= 20.19, 与 CI 中使用的版本保持一致
FROM node:22-alpine AS BUILD
RUN npm config set registry https://registry.npmmirror.com \
    && npm install -g pnpm@11

# 依赖清单单独拷贝, 依赖没变时可以复用这一层缓存
# pnpm-workspace.yaml 里有 allowBuilds 配置, 不拷会报 ERR_PNPM_IGNORED_BUILDS
WORKDIR /app/front
COPY gin-blog-front/package.json gin-blog-front/pnpm-lock.yaml gin-blog-front/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY gin-blog-front .
RUN pnpm build

WORKDIR /app/admin
COPY gin-blog-admin/package.json gin-blog-admin/pnpm-lock.yaml gin-blog-admin/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY gin-blog-admin .
RUN pnpm build

# 阶段二: 将静态资源部署到 Nginx
FROM nginx:1.24.0-alpine

# 从第一个阶段拷贝构建好的静态资源到容器
COPY --from=BUILD /app/front/dist /usr/share/nginx/html
COPY --from=BUILD /app/admin/dist /usr/share/nginx/html/admin

# 将 Nginx 配置文件模板拷到容器中, 并执行脚本填充环境变量
COPY deploy/build/web/default.conf.template /etc/nginx/conf.d/default.conf.template
COPY deploy/build/web/default.conf.ssl.template /etc/nginx/conf.d/default.conf.ssl.template
COPY deploy/build/web/run.sh /docker-entrypoint.sh
RUN chmod a+x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]

CMD [ "nginx", "-g", "daemon off;" ]

EXPOSE 80
EXPOSE 443