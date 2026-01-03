package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	// 服务器配置
	Port string

	// 数据库配置
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT 配置
	JWTSecret string

	// 上传配置
	UploadPath string
	UploadURL  string

	// 管理员配置
	AdminUsername string
	AdminPassword string
}

// Load 加载配置
func Load() *Config {
	// 尝试加载 .env 文件
	godotenv.Load()

	return &Config{
		Port:          getEnv("PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "dualtab"),
		DBPassword:    getEnv("DB_PASSWORD", "dualtab123"),
		DBName:        getEnv("DB_NAME", "dualtab"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		UploadPath:    getEnv("UPLOAD_PATH", "./uploads"),
		UploadURL:     getEnv("UPLOAD_URL", "http://localhost:8080/uploads"),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
