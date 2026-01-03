# 生产环境部署指南

## 🎯 运行时环境变量方案

采用 **`window.__ENV__` 运行时注入** 的标准模式,完全支持预构建镜像。

### 原理

通过服务端组件在 HTML 中注入全局变量:

```html
<script>
  window.__ENV__ = { API_URL: "http://your-server:8080" }
</script>
```

前端代码读取 `window.__ENV__.API_URL` 作为 API 地址。

**优势**:
- ✅ 真正的运行时配置,无需重新构建
- ✅ 直接使用预构建镜像
- ✅ 修改配置只需重启容器
- ✅ 不依赖 Nginx,灵活适配各种部署环境

---

## 快速部署

### 1. 准备环境

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | sh
```

### 2. 下载部署文件

```bash
mkdir dualtab-backend && cd dualtab-backend

# 下载配置文件
curl -O https://raw.githubusercontent.com/EmccK/DualTabBackend/main/docker-compose.prod.yml
curl -O https://raw.githubusercontent.com/EmccK/DualTabBackend/main/.env.prod.example

# 创建数据目录
mkdir -p data/postgres uploads
```

### 3. 配置环境变量

```bash
cp .env.prod.example .env.prod
vim .env.prod
```

修改以下配置:

```bash
# 数据库密码
POSTGRES_PASSWORD=your-strong-password-here

# JWT 密钥 (生成: openssl rand -base64 32)
JWT_SECRET=your-random-jwt-secret-key-here

# 管理员密码
ADMIN_PASSWORD=your-admin-password-here

# ⚠️ 前端 API 地址 (改为你的服务器 IP 或域名)
API_URL=http://123.45.67.89:8080

# 上传文件 URL
UPLOAD_URL=http://123.45.67.89:8080/uploads
```

### 4. 启动服务

**✅ 直接启动,无需 --build!**

```bash
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

### 5. 验证

访问: `http://你的IP:3000`

打开浏览器开发者工具:
1. Console 输入 `window.__ENV__` 查看注入的配置
2. Network 标签查看请求地址

应该看到:
- ✅ `window.__ENV__.API_URL` 为你配置的地址
- ✅ 请求发往 `http://你的IP:8080/admin/...`

---

## 配置说明

### 环境变量优先级

```
运行时 API_URL > 构建时 NEXT_PUBLIC_API_URL > 默认 localhost:8080
```

### 不同部署场景

#### 场景 1: 直接暴露端口

```bash
# .env.prod
API_URL=http://123.45.67.89:8080
BACKEND_PORT=8080
ADMIN_PORT=3000
```

访问: `http://123.45.67.89:3000`

#### 场景 2: 使用自己的 Nginx

```bash
# .env.prod
API_URL=https://yourdomain.com/api
BACKEND_PORT=8080
ADMIN_PORT=3000
```

Nginx 配置示例:
```nginx
location /api/ {
    proxy_pass http://localhost:8080/;
}
location / {
    proxy_pass http://localhost:3000;
}
```

访问: `https://yourdomain.com`

#### 场景 3: 前后端同域名

```bash
# .env.prod
API_URL=https://yourdomain.com
```

Nginx 配置:
```nginx
location /admin { proxy_pass http://localhost:8080; }
location / { proxy_pass http://localhost:3000; }
```

---

## 修改配置

如果修改了 `API_URL`,只需重启容器:

```bash
# 编辑配置
vim .env.prod

# 重启(无需 --build)
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d

# 或单独重启前端
docker-compose -f docker-compose.prod.yml restart admin
```

**1-2 秒即可生效!**

---

## 更新镜像

```bash
# 拉取最新镜像
docker-compose -f docker-compose.prod.yml pull

# 重启
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d

# 清理旧镜像
docker image prune -f
```

---

## 备份与恢复

### 备份

```bash
# 备份数据库
docker exec dualtab-db pg_dump -U dualtab dualtab > backup-$(date +%Y%m%d).sql

# 备份数据目录
tar -czf backup-$(date +%Y%m%d).tar.gz data/ uploads/
```

### 恢复

```bash
# 恢复数据库
docker exec -i dualtab-db psql -U dualtab dualtab < backup-20260103.sql
```

---

## 常见问题

### Q1: 还需要配置 NEXT_PUBLIC_API_URL 吗?

❌ 不需要。现在使用 `API_URL` 运行时环境变量。

### Q2: 还需要 --build 吗?

❌ 不需要。直接 `docker-compose up -d` 即可。

### Q3: 修改 API_URL 后需要重新构建吗?

❌ 不需要。只需重启容器即可,1-2秒生效。

### Q4: 可以用自己的 Nginx 吗?

✅ 可以。随意配置,只需在 `.env.prod` 中设置对应的 `API_URL`。

### Q5: 如何验证配置是否生效?

浏览器 Console 输入:
```javascript
window.__ENV__
```

应该看到:
```json
{ "API_URL": "http://你配置的地址:8080" }
```

---

## 参考资料

- [Next.js Environment Variables](https://nextjs.org/docs/pages/guides/environment-variables)
- [Runtime Environment Variables in Next.js Docker](https://dev.to/nemanjam/runtime-environment-variables-in-nextjs-build-reusable-docker-images-ho)
- [Next.js Runtime Config Discussion](https://github.com/vercel/next.js/discussions/44628)
