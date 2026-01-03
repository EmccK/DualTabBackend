package api

import (
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// CategoryHandler 对外分类 API 处理器
type CategoryHandler struct {
	repo *repository.CategoryRepo
}

// NewCategoryHandler 创建对外分类 API 处理器
func NewCategoryHandler(repo *repository.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// GetList 获取分类列表
// GET /categories
func (h *CategoryHandler) GetList(c *gin.Context) {
	categories, err := h.repo.FindAll(true) // 只返回激活的
	if err != nil {
		response.Success(c, gin.H{"list": []interface{}{}})
		return
	}

	// 转换为前端需要的格式
	list := make([]gin.H, len(categories))
	for i, cat := range categories {
		list[i] = gin.H{
			"id":      cat.ID,
			"name":    cat.Name,
			"name_en": cat.NameEn,
		}
	}

	response.Success(c, gin.H{"list": list})
}
