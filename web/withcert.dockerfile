# ──────────────────────────────────────────────────────────
#  1. 选择镜像
# ──────────────────────────────────────────────────────────
FROM nginx:alpine

# ──────────────────────────────────────────────────────────
#  2. 将 Cloudflare Origin 证书拷贝进镜像
#    假设你在项目根目录放置：
#      certs/
#        ├─ origin.pem   # Cloudflare Origin Certificate
#        └─ origin.key   # 你生成的私钥
# ──────────────────────────────────────────────────────────
RUN mkdir -p /etc/nginx/ssl
COPY certs/origin.pem /etc/nginx/ssl/origin.pem
COPY certs/origin.key /etc/nginx/ssl/origin.key

# 权限收紧（可选，但推荐）
RUN chmod 600 /etc/nginx/ssl/origin.key && \
    chmod 644 /etc/nginx/ssl/origin.pem

# ──────────────────────────────────────────────────────────
#  3. 覆盖 Nginx 配置
# ──────────────────────────────────────────────────────────
RUN rm /etc/nginx/conf.d/default.conf
RUN cat <<'EOF' > /etc/nginx/conf.d/default.conf
# --- HTTP 端口：全部 301 跳转到 HTTPS ---------------------------------
server {
    listen 80;
    server_name example.invalid;           # ← 请改成你的域名
    return 301 https://$host$request_uri;
}

# --- HTTPS 端口 --------------------------------------------------------
server {
    listen 443 ssl http2;
    server_name example.invalid;           # ← 同上

    # SSL 证书
    ssl_certificate     /etc/nginx/ssl/origin.pem;
    ssl_certificate_key /etc/nginx/ssl/origin.key;

    # 基础安全 / 性能参数
    ssl_protocols              TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers  off;
    ssl_session_cache          shared:SSL:10m;
    ssl_session_timeout        1h;

    # 静态资源
    root /usr/share/nginx/html;
    index index.html index.htm;

    # SPA History 模式回退
    location / {
        try_files $uri $uri/ /index.html;
    }

    # WebSocket：/maasst/logs-ws
    location /maasst/logs-ws {
        proxy_pass http://127.0.0.1:5244;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 86400;
    }

    # 其他后端接口：/maasst/*
    location /maasst/ {
        proxy_pass http://127.0.0.1:5244;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
EOF

# ──────────────────────────────────────────────────────────
# 4. 拷贝前端静态文件
# ──────────────────────────────────────────────────────────
COPY dist/ /usr/share/nginx/html/

# ──────────────────────────────────────────────────────────
# 5. 暴露端口 & 启动
# ──────────────────────────────────────────────────────────
EXPOSE 80 443
CMD ["nginx", "-g", "daemon off;"]