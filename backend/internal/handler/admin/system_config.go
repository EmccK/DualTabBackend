package admin

import (
	"dualtab-backend/internal/model"
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// SystemConfigHandler 系统配置管理处理器
type SystemConfigHandler struct {
	repo *repository.SystemConfigRepo
}

// NewSystemConfigHandler 创建系统配置管理处理器
func NewSystemConfigHandler(repo *repository.SystemConfigRepo) *SystemConfigHandler {
	return &SystemConfigHandler{repo: repo}
}

// List 获取所有配置
func (h *SystemConfigHandler) List(c *gin.Context) {
	configs, err := h.repo.FindAll()
	if err != nil {
		response.InternalError(c, "获取配置列表失败")
		return
	}
	response.Success(c, gin.H{"list": configs})
}

// Get 获取单个配置
func (h *SystemConfigHandler) Get(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "请提供配置 Key")
		return
	}

	config, err := h.repo.FindByKey(key)
	if err != nil {
		response.NotFound(c, "配置不存在")
		return
	}

	response.Success(c, config)
}

// SetConfigRequest 设置配置请求
type SetConfigRequest struct {
	Key    string `json:"key" binding:"required"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

// Set 设置配置
func (h *SystemConfigHandler) Set(c *gin.Context) {
	var req SetConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请提供配置 Key")
		return
	}

	if err := h.repo.SetValue(req.Key, req.Value, req.Remark); err != nil {
		response.InternalError(c, "设置配置失败")
		return
	}

	response.Success(c, gin.H{
		"key":   req.Key,
		"value": req.Value,
	})
}

// BatchSetRequest 批量设置配置请求
type BatchSetRequest struct {
	Configs []SetConfigRequest `json:"configs" binding:"required"`
}

// BatchSet 批量设置配置
func (h *SystemConfigHandler) BatchSet(c *gin.Context) {
	var req BatchSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请提供配置列表")
		return
	}

	// 转换为仓库层需要的格式
	configs := make([]struct {
		Key    string
		Value  string
		Remark string
	}, len(req.Configs))

	for i, cfg := range req.Configs {
		configs[i].Key = cfg.Key
		configs[i].Value = cfg.Value
		configs[i].Remark = cfg.Remark
	}

	// 使用事务批量设置
	if err := h.repo.BatchSetValue(configs); err != nil {
		response.InternalError(c, "批量设置配置失败")
		return
	}

	response.Success(c, nil)
}

// Delete 删除配置
func (h *SystemConfigHandler) Delete(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "请提供配置 Key")
		return
	}

	if err := h.repo.Delete(key); err != nil {
		response.InternalError(c, "删除配置失败")
		return
	}

	response.Success(c, nil)
}

// GetConfigKeys 获取所有可用的配置 Key 说明
func (h *SystemConfigHandler) GetConfigKeys(c *gin.Context) {
	keys := []gin.H{
		{
			"key":         model.ConfigKeyWeatherAPIKey,
			"description": "天气 API Key（和风天气或 OpenWeather）",
			"example":     "your-api-key",
		},
		{
			"key":         model.ConfigKeyWeatherAPIType,
			"description": "天气 API 类型",
			"example":     "qweather 或 openweather",
		},
		{
			"key":         model.ConfigKeySearchSuggestOn,
			"description": "是否启用搜索建议代理",
			"example":     "true 或 false",
		},
		{
			"key":         model.ConfigKeyBingWallpaperOn,
			"description": "是否启用 Bing 每日壁纸（当没有自定义壁纸时使用）",
			"example":     "true 或 false",
		},
	}

	response.Success(c, gin.H{"keys": keys})
}
