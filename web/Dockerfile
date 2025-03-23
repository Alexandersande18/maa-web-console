# 使用官方的 Nginx alpine 镜像
FROM nginx:alpine

# 删除 Nginx 默认的配置文件，避免冲突
RUN rm /etc/nginx/conf.d/default.conf

# 使用 here document 方式创建自定义的 Nginx 配置文件
RUN cat <<'EOF' > /etc/nginx/conf.d/default.conf
server {
    listen 80;
    server_name example.invalid;  # 请根据实际情况替换为你的域名或IP

    # 配置静态资源的根目录及默认首页文件
    root /usr/share/nginx/html;
    index index.html index.htm;

    # 前端路由和静态文件支持：当请求的资源不存在时，回退到 index.html
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 针对 websocket 接口 /maasst/logs-ws
    location /maasst/logs-ws {
        proxy_pass http://127.0.0.1:5244;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 86400;
    }

    # 针对其他 GoFrame 后端接口 /maasst/*
    location /maasst/ {
        proxy_pass http://127.0.0.1:5244;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
EOF

# 将构建后的 Vue 静态文件拷贝到 Nginx 默认的静态文件目录
COPY dist/ /usr/share/nginx/html/

# 暴露 Nginx 的默认端口
EXPOSE 80

# 启动 Nginx
CMD ["nginx", "-g", "daemon off;"]