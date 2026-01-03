# 生产环境部署指南

## 快速部署

### 1. 准备环境

确保服务器已安装 Docker 和 Docker Compose：
```bash
docker --version
docker-compose --version
```

### 2. 下载部署文件

```bash
# 创建项目目录
mkdir dualtab-backend && cd dualtab-backend

# 下载部署配置文件
curl -O https://raw.githubusercontent.com/EmccK/DualTabBackend/main/docker-compose.prod.yml
curl -O https://raw.githubusercontent.com/EmccK/DualTabBackend/main/.env.prod.example

# 重命名环境变量文件
mv .env.prod.example .env

# 创建数据目录
mkdir -p data/postgres uploads
```

### 3. 配置环境变量

编辑 `.env` 文件，修改以下关键配置：

```bash
# 修改数据库密码
POSTGRES_PASSWORD=your-strong-password-here

# 修改 JWT 密钥（至少 32 位随机字符串）
JWT_SECRET=your-random-jwt-secret-key-here

# 修改管理员密码
ADMIN_PASSWORD=your-admin-password-here

# 修改域名（如果有域名）
UPLOAD_URL=http://your-domain.com:8080/uploads
NEXT_PUBLIC_API_URL=http://your-domain.com:8080
```

生成随机 JWT 密钥：
```bash
openssl rand -base64 32
```

### 4. 启动服务

```bash
# 拉取最新镜像
docker-compose -f docker-compose.prod.yml pull

# 启动服务
docker-compose -f docker-compose.prod.yml up -d

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

### 5. 访问服务

- 后端 API: http://your-ip:8080
- 管理后台: http://your-ip:3000

默认管理员账号：
- 用户名: `admin` (可在 .env 中修改)
- 密码: 你在 .env 中设置的密码

## 更新镜像

```bash
# 拉取最新镜像
docker-compose -f docker-compose.prod.yml pull

# 重启服务
docker-compose -f docker-compose.prod.yml up -d

# 清理旧镜像
docker image prune -f
```

## 备份数据

### 备份数据库

```bash
# 导出数据库
docker exec dualtab-db pg_dump -U dualtab dualtab > backup-$(date +%Y%m%d).sql

# 或直接备份数据目录
tar -czf data-backup-$(date +%Y%m%d).tar.gz data/

# 恢复数据库
docker exec -i dualtab-db psql -U dualtab dualtab < backup-20260103.sql
```

### 备份上传文件

```bash
# 备份上传目录
tar -czf uploads-backup-$(date +%Y%m%d).tar.gz uploads/
```

## 使用 Nginx 反向代理（推荐）

创建 `/etc/nginx/sites-available/dualtab` 配置：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 后端 API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 上传文件
    location /uploads {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
    }

    # 管理后台
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用配置并重启 Nginx：
```bash
sudo ln -s /etc/nginx/sites-available/dualtab /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 常用命令

```bash
# 查看服务状态
docker-compose -f docker-compose.prod.yml ps

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f [service_name]

# 重启服务
docker-compose -f docker-compose.prod.yml restart

# 停止服务
docker-compose -f docker-compose.prod.yml down

# 完全清理（警告：会删除所有数据，请先备份）
# 注意：使用路径方式时，需要手动删除 data 目录
docker-compose -f docker-compose.prod.yml down
rm -rf data/ uploads/
```

## 安全建议

1. 修改默认密码和密钥
2. 使用 HTTPS（配置 SSL 证书）
3. 配置防火墙，只开放必要端口
4. 定期备份数据库和上传文件
5. 定期更新镜像到最新版本
6. 使用强密码和随机 JWT 密钥
