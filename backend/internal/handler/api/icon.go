package api

import (
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IconHandler 对外图标 API 处理器
type IconHandler struct {
	repo *repository.IconRepo
}

// NewIconHandler 创建对外图标 API 处理器
func NewIconHandler(repo *repository.IconRepo) *IconHandler {
	return &IconHandler{repo: repo}
}

// GetList 获取推荐书签列表（兼容 MonkNow 格式）
// GET /icon/list?cate_id=24&keyword=&size=20
func (h *IconHandler) GetList(c *gin.Context) {
	cateID, _ := strconv.ParseUint(c.Query("cate_id"), 10, 32)
	keyword := c.Query("keyword")
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	// 只查询激活的图标
	isActive := true
	query := repository.IconQuery{
		CategoryID: uint(cateID),
		Keyword:    keyword,
		IsActive:   &isActive,
		Size:       size,
	}

	icons, _, err := h.repo.FindAll(query)
	if err != nil {
		response.Success(c, gin.H{"list": []interface{}{}})
		return
	}

	// 转换为 MonkNow 格式
	list := make([]gin.H, len(icons))
	for i, icon := range icons {
		list[i] = gin.H{
			"udId":        icon.ID,
			"title":       icon.Title,
			"description": icon.Description,
			"url":         icon.URL,
			"imgUrl":      icon.ImgURL,
			"bgColor":     icon.BgColor,
			"mimeType":    icon.MimeType,
		}
	}

	response.Success(c, gin.H{"list": list})
}

// GetByURL 根据 URL 获取图标
// GET /icon/byurl?url=https://example.com
func (h *IconHandler) GetByURL(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		response.BadRequest(c, "请提供 URL")
		return
	}

	icon, err := h.repo.FindByURL(url)
	if err != nil {
		// 未找到时返回空数据
		response.Success(c, nil)
		return
	}

	response.Success(c, gin.H{
		"udId":        icon.ID,
		"title":       icon.Title,
		"description": icon.Description,
		"url":         icon.URL,
		"imgUrl":      icon.ImgURL,
		"bgColor":     icon.BgColor,
		"mimeType":    icon.MimeType,
	})
}
