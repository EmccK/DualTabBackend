# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装时区数据
RUN apk --no-cache add tzdata

COPY --from=builder /app/main .

# 创建上传目录
RUN mkdir -p /app/uploads/icons

EXPOSE 8080

CMD ["./main"]
