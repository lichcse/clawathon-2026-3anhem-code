FROM node:20-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

COPY <<EOF /etc/nginx/conf.d/default.conf
server {
  listen 5173;
  location /api/ {
    proxy_pass http://backend:8080;
    proxy_set_header Host \$host;
    proxy_set_header X-Real-IP \$remote_addr;
  }
  location /ws {
    proxy_pass http://backend:8080;
    proxy_http_version 1.1;
    proxy_set_header Upgrade \$http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host \$host;
  }
  location / {
    root /usr/share/nginx/html;
    index index.html index.htm;
    try_files \$uri \$uri/ /index.html;
  }
}
EOF

EXPOSE 5173
CMD ["nginx", "-g", "daemon off;"]
