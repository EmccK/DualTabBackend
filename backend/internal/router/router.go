package router

import (
	"dualtab-backend/config"
	adminHandler "dualtab-backend/internal/handler/admin"
	apiHandler "dualtab-backend/internal/handler/api"
	"dualtab-backend/internal/middleware"
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/upload"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup 配置路由
func Setup(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// 创建上传器
	uploader := upload.NewUploader(cfg.UploadPath, cfg.UploadURL)

	// 创建仓库
	adminUserRepo := repository.NewAdminUserRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	iconRepo := repository.NewIconRepo(db)
	searchEngineRepo := repository.NewSearchEngineRepo(db)

	// 创建管理后台处理器
	authHandler := adminHandler.NewAuthHandler(adminUserRepo, cfg.JWTSecret)
	adminCategoryHandler := adminHandler.NewCategoryHandler(categoryRepo)
	adminIconHandler := adminHandler.NewIconHandler(iconRepo, uploader)
	adminSearchEngineHandler := adminHandler.NewSearchEngineHandler(searchEngineRepo)

	// 创建对外 API 处理器
	apiIconHandler := apiHandler.NewIconHandler(iconRepo)
	apiSearchEngineHandler := apiHandler.NewSearchEngineHandler(searchEngineRepo)
	apiCategoryHandler := apiHandler.NewCategoryHandler(categoryRepo)

	// 中间件
	r.Use(middleware.CORS())

	// 静态文件服务
	r.Static("/uploads", cfg.UploadPath)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// ========== 对外 API（兼容 MonkNow）==========
	r.GET("/icon/list", apiIconHandler.GetList)
	r.GET("/icon/byurl", apiIconHandler.GetByURL)
	r.GET("/search-engines", apiSearchEngineHandler.GetList)
	r.GET("/categories", apiCategoryHandler.GetList)

	// ========== 管理后台 API ==========
	adminGroup := r.Group("/admin")
	{
		// 认证（无需登录）
		adminGroup.POST("/auth/login", authHandler.Login)
		adminGroup.POST("/auth/register", authHandler.CreateAdmin) // 仅用于初始化

		// 需要认证的路由
		authGroup := adminGroup.Group("")
		authGroup.Use(middleware.AdminAuth(cfg.JWTSecret))
		{
			// 当前用户
			authGroup.GET("/auth/me", authHandler.GetCurrentUser)
			authGroup.PUT("/auth/password", authHandler.ChangePassword)

			// 分类管理
			authGroup.GET("/categories", adminCategoryHandler.List)
			authGroup.GET("/categories/:id", adminCategoryHandler.Get)
			authGroup.POST("/categories", adminCategoryHandler.Create)
			authGroup.PUT("/categories/:id", adminCategoryHandler.Update)
			authGroup.DELETE("/categories/:id", adminCategoryHandler.Delete)

			// 图标管理
			authGroup.GET("/icons", adminIconHandler.List)
			authGroup.GET("/icons/:id", adminIconHandler.Get)
			authGroup.POST("/icons", adminIconHandler.Create)
			authGroup.PUT("/icons/:id", adminIconHandler.Update)
			authGroup.DELETE("/icons/:id", adminIconHandler.Delete)

			// 搜索引擎管理
			authGroup.GET("/search-engines", adminSearchEngineHandler.List)
			authGroup.GET("/search-engines/:id", adminSearchEngineHandler.Get)
			authGroup.POST("/search-engines", adminSearchEngineHandler.Create)
			authGroup.PUT("/search-engines/:id", adminSearchEngineHandler.Update)
			authGroup.DELETE("/search-engines/:id", adminSearchEngineHandler.Delete)

			// 文件上传
			authGroup.POST("/upload/icon", adminIconHandler.UploadIcon)
		}
	}
}
