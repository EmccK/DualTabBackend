package main

import (
	"fmt"
	"log"

	"dualtab-backend/config"
	"dualtab-backend/internal/model"
	"dualtab-backend/internal/router"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 连接数据库
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&model.AdminUser{},
		&model.Category{},
		&model.Icon{},
		&model.SearchEngine{},
		&model.Wallpaper{},
		&model.SystemConfig{},
		&model.User{},
		&model.UserData{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化预设数据
	initPresetData(db, cfg)

	// 创建 Gin 引擎
	r := gin.Default()

	// 配置路由
	router.Setup(r, db, cfg)

	// 启动服务
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("服务启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}

// initPresetData 初始化预设数据
func initPresetData(db *gorm.DB, cfg *config.Config) {
	// 预设分类
	categories := []model.Category{
		{ID: 24, Name: "热门", NameEn: "hot", SortOrder: 1, IsActive: true},
		{ID: 9, Name: "购物", NameEn: "shopping", SortOrder: 2, IsActive: true},
		{ID: 10, Name: "社交", NameEn: "social", SortOrder: 3, IsActive: true},
		{ID: 26, Name: "娱乐", NameEn: "entertainment", SortOrder: 4, IsActive: true},
		{ID: 11, Name: "新闻与阅读", NameEn: "news", SortOrder: 5, IsActive: true},
		{ID: 14, Name: "效率", NameEn: "efficiency", SortOrder: 6, IsActive: true},
		{ID: 25, Name: "内置App", NameEn: "builtin", SortOrder: 7, IsActive: true},
		{ID: 15, Name: "图片", NameEn: "image", SortOrder: 8, IsActive: true},
		{ID: 16, Name: "生活方式", NameEn: "lifestyle", SortOrder: 9, IsActive: true},
		{ID: 17, Name: "旅行", NameEn: "travel", SortOrder: 10, IsActive: true},
		{ID: 18, Name: "科技与教育", NameEn: "tech", SortOrder: 11, IsActive: true},
		{ID: 19, Name: "金融", NameEn: "finance", SortOrder: 12, IsActive: true},
	}

	for _, cat := range categories {
		db.FirstOrCreate(&cat, model.Category{ID: cat.ID})
	}

	// 预设搜索引擎
	searchEngines := []model.SearchEngine{
		{UUID: "e58b5a00-74fe-4319-af0a-d4999565dd71", Name: "Google", URL: "https://www.google.com/search?q=", IconURL: "https://static.monknow.com/newtab/searcher/e58b5a00-74fe-4319-af0a-d4999565dd71.svg", SortOrder: 1, IsActive: true},
		{UUID: "0eb43a90-b4c7-43ce-9c73-ab110945f47d", Name: "百度", URL: "https://www.baidu.com/s?wd=", IconURL: "https://static.monknow.com/newtab/searcher/0eb43a90-b4c7-43ce-9c73-ab110945f47d.svg", SortOrder: 2, IsActive: true},
		{UUID: "ceb6c985-d09c-4fdc-b0ea-b304f1ee0f2d", Name: "Bing", URL: "https://www.bing.com/search?q=", IconURL: "https://static.monknow.com/newtab/searcher/ceb6c985-d09c-4fdc-b0ea-b304f1ee0f2d.svg", SortOrder: 3, IsActive: true},
		{UUID: "2a5e69d9-bf13-4188-8da2-004551a913a0", Name: "Yahoo", URL: "https://search.yahoo.com/search?p=", IconURL: "https://static.monknow.com/newtab/searcher/2a5e69d9-bf13-4188-8da2-004551a913a0.svg", SortOrder: 4, IsActive: true},
		{UUID: "118f7463-4411-4856-873f-2851faa3b543", Name: "Yandex", URL: "https://yandex.ru/search/?text=", IconURL: "https://static.monknow.com/newtab/searcher/118f7463-4411-4856-873f-2851faa3b543.svg", SortOrder: 5, IsActive: true},
		{UUID: "259d8e2b-340e-4690-8046-88a0b130cbd0", Name: "DuckDuckGo", URL: "https://duckduckgo.com/?q=", IconURL: "https://static.monknow.com/newtab/searcher/259d8e2b-340e-4690-8046-88a0b130cbd0.svg", SortOrder: 6, IsActive: true},
	}

	for _, engine := range searchEngines {
		db.FirstOrCreate(&engine, model.SearchEngine{UUID: engine.UUID})
	}

	// 初始化管理员账号（从环境变量读取）
	var adminCount int64
	db.Model(&model.AdminUser{}).Count(&adminCount)
	if adminCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("生成密码哈希失败: %v", err)
			return
		}
		admin := model.AdminUser{
			Username:     cfg.AdminUsername,
			PasswordHash: string(hash),
		}
		db.Create(&admin)
		log.Printf("已创建管理员账号: %s", cfg.AdminUsername)
	}
}
