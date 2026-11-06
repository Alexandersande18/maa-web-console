# Certificates must be mounted at runtime, never copied into image layers.
FROM nginx:alpine
COPY nginx.tls.conf.example /etc/nginx/conf.d/default.conf
COPY dist/ /usr/share/nginx/html/
EXPOSE 80 443
CMD ["nginx", "-g", "daemon off;"]
