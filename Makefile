.PHONY: dev build run docker-up docker-down clean

# 开发模式运行后端
dev:
	cd backend && go run main.go

# 编译后端
build:
	cd backend && go build -o ../bin/backend main.go

# 运行编译后的后端
run:
	./bin/backend

# 启动 Docker 容器
docker-up:
	docker-compose up -d

# 停止 Docker 容器
docker-down:
	docker-compose down

# 重新构建并启动
docker-rebuild:
	docker-compose up -d --build

# 查看日志
logs:
	docker-compose logs -f

# 清理
clean:
	rm -rf bin/
	docker-compose down -v

# 安装后端依赖
deps:
	cd backend && go mod tidy

# 初始化管理前端
init-admin:
	cd admin && npx create-next-app@latest . --typescript --tailwind --eslint --app --src-dir --import-alias "@/*"
