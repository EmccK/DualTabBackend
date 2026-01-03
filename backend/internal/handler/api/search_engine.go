package api

import (
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// SearchEngineHandler 对外搜索引擎 API 处理器
type SearchEngineHandler struct {
	repo *repository.SearchEngineRepo
}

// NewSearchEngineHandler 创建对外搜索引擎 API 处理器
func NewSearchEngineHandler(repo *repository.SearchEngineRepo) *SearchEngineHandler {
	return &SearchEngineHandler{repo: repo}
}

// GetList 获取搜索引擎列表
// GET /search-engines
func (h *SearchEngineHandler) GetList(c *gin.Context) {
	engines, err := h.repo.FindAll(true) // 只返回激活的
	if err != nil {
		response.Success(c, gin.H{"list": []interface{}{}})
		return
	}

	// 转换为前端需要的格式
	list := make([]gin.H, len(engines))
	for i, engine := range engines {
		list[i] = gin.H{
			"id":   engine.UUID,
			"name": engine.Name,
			"url":  engine.URL,
			"icon": engine.IconURL,
		}
	}

	response.Success(c, gin.H{"list": list})
}
