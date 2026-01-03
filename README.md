# DualTab 后台管理系统

DualTab Chrome 扩展的后端服务和管理后台。

## 技术栈

- **后端**: Go + Gin + GORM
- **数据库**: PostgreSQL
- **管理前端**: Next.js + shadcn/ui
- **部署**: Docker

## 快速开始

### 本地开发环境

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

默认服务地址:
- 后端 API: http://localhost:8080
- 管理后台: http://localhost:3000

### 生产环境部署

使用 **运行时环境变量** 方案,无需重新构建镜像。

详细部署步骤: [DEPLOY.md](./DEPLOY.md)

快速部署:

```bash
# 1. 下载配置文件
curl -O https://raw.githubusercontent.com/EmccK/DualTabBackend/main/docker-compose.prod.yml
curl -O https://raw.githubusercontent.com/EmccK/DualTabBackend/main/.env.prod.example

# 2. 配置环境变量
cp .env.prod.example .env.prod
vim .env.prod  # 修改 API_URL、密码等

# 3. 启动服务(无需 --build)
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d
```

访问: `http://你的IP:3000`

### 本地开发

**后端**

```bash
cd backend

# 安装依赖
go mod tidy

# 复制环境变量
cp .env.example .env

# 启动 PostgreSQL（需要先安装）
# 或使用 Docker: docker-compose up -d postgres

# 运行
go run main.go
```

**管理前端**

```bash
cd admin

# 安装依赖
npm install

# 运行
npm run dev
```

## 默认账号

- 用户名: `admin`
- 密码: `admin123`

## API 文档

### 对外 API（供 DualTab 扩展调用）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/icon/list` | 获取推荐书签列表 |
| GET | `/icon/byurl` | 根据 URL 获取图标 |
| GET | `/search-engines` | 获取搜索引擎列表 |
| GET | `/categories` | 获取分类列表 |
| GET | `/wallpaper/random` | 获取随机壁纸 |
| GET | `/wallpaper/list` | 获取壁纸列表 |
| GET | `/weather/locations` | 搜索城市位置 |
| GET | `/weather` | 获取天气信息 |
| GET | `/proxy/search-suggest` | 搜索建议代理 |

### 用户 API（数据同步）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/user/login` | 用户登录 |
| POST | `/user/register` | 用户注册 |
| GET | `/home/captcha` | 获取验证码 |
| GET | `/user/data/info` | 获取用户数据（需 secret 请求头） |
| PUT | `/user/data/update` | 更新用户数据（需 secret 请求头） |
| PUT | `/user/changename` | 修改昵称（需 secret 请求头） |
| PUT | `/user/changeavatar` | 修改头像（需 secret 请求头） |
| PUT | `/user/changepwd` | 修改密码（需 secret 请求头） |
| POST | `/upload/image` | 上传图片 |

### 管理后台 API

需要在请求头中携带 `Authorization: Bearer <token>`

**认证**
- POST `/admin/auth/login` - 登录
- GET `/admin/auth/me` - 获取当前用户

**图标管理**
- GET `/admin/icons` - 图标列表
- POST `/admin/icons` - 创建图标
- PUT `/admin/icons/:id` - 更新图标
- DELETE `/admin/icons/:id` - 删除图标

**分类管理**
- GET `/admin/categories` - 分类列表
- POST `/admin/categories` - 创建分类
- PUT `/admin/categories/:id` - 更新分类
- DELETE `/admin/categories/:id` - 删除分类

**搜索引擎管理**
- GET `/admin/search-engines` - 搜索引擎列表
- POST `/admin/search-engines` - 创建搜索引擎
- PUT `/admin/search-engines/:id` - 更新搜索引擎
- DELETE `/admin/search-engines/:id` - 删除搜索引擎

**壁纸管理**
- GET `/admin/wallpapers` - 壁纸列表
- POST `/admin/wallpapers` - 创建壁纸
- PUT `/admin/wallpapers/:id` - 更新壁纸
- DELETE `/admin/wallpapers/:id` - 删除壁纸

**系统配置**
- GET `/admin/configs` - 配置列表
- GET `/admin/configs/keys` - 获取可用配置项说明
- POST `/admin/configs` - 设置配置
- POST `/admin/configs/batch` - 批量设置配置

**文件上传**
- POST `/admin/upload/icon` - 上传图标图片
- POST `/admin/upload/wallpaper` - 上传壁纸图片

## 系统配置项

| Key | 说明 | 示例值 |
|-----|------|--------|
| `weather_api_key` | 天气 API Key | your-api-key |
| `weather_api_type` | 天气 API 类型 | qweather / openweather |
| `search_suggest_on` | 启用搜索建议代理 | true / false |
| `bing_wallpaper_on` | 启用 Bing 每日壁纸 | true / false |

## 目录结构

```
.
├── backend/                 # Go 后端
│   ├── config/             # 配置
│   ├── internal/
│   │   ├── handler/        # HTTP 处理器
│   │   ├── middleware/     # 中间件
│   │   ├── model/          # 数据模型
│   │   ├── repository/     # 数据访问层
│   │   ├── service/        # 业务服务
│   │   └── router/         # 路由
│   ├── pkg/                # 公共包
│   └── migrations/         # 数据库迁移
├── admin/                   # Next.js 管理前端
│   └── src/
│       ├── app/            # 页面
│       ├── components/     # 组件
│       └── lib/            # 工具函数
├── uploads/                 # 上传文件
└── docker-compose.yml
```

## 环境变量

### 后端

| 变量 | 说明 | 默认值 |
|------|------|--------|
| PORT | 服务端口 | 8080 |
| DB_HOST | 数据库地址 | localhost |
| DB_PORT | 数据库端口 | 5432 |
| DB_USER | 数据库用户 | dualtab |
| DB_PASSWORD | 数据库密码 | dualtab123 |
| DB_NAME | 数据库名 | dualtab |
| JWT_SECRET | JWT 密钥 | - |
| UPLOAD_PATH | 上传目录 | ./uploads |
| UPLOAD_URL | 上传文件访问 URL | http://localhost:18080/uploads |

### 管理前端

| 变量 | 说明 | 默认值 |
|------|------|--------|
| API_URL | 后端 API 地址(运行时配置) | http://localhost:8080 |

**注意**: 使用 `API_URL` 运行时环境变量,修改后只需重启容器即可生效。
