package api

import (
	"dualtab-backend/internal/repository"
	"dualtab-backend/internal/service"
	"dualtab-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IconHandler 对外图标 API 处理器
type IconHandler struct {
	repo           *repository.IconRepo
	faviconService *service.FaviconService
}

// NewIconHandler 创建对外图标 API 处理器
func NewIconHandler(repo *repository.IconRepo, faviconService *service.FaviconService) *IconHandler {
	return &IconHandler{
		repo:           repo,
		faviconService: faviconService,
	}
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

	// 1. 先尝试从icons表查找（推荐图标）
	icon, err := h.repo.FindByURL(url)
	if err == nil {
		// 找到推荐图标，直接返回
		response.Success(c, gin.H{
			"udId":        icon.ID,
			"title":       icon.Title,
			"description": icon.Description,
			"url":         icon.URL,
			"imgUrl":      icon.ImgURL,
			"bgColor":     icon.BgColor,
			"mimeType":    icon.MimeType,
		})
		return
	}

	// 2. 推荐图标未找到，使用FaviconService获取
	// FaviconService会先从缓存查找，找不到再调用API
	favicon, err := h.faviconService.GetFavicon(url)
	if err != nil {
		// 获取失败，返回空数据（保持与前端兼容）
		response.Success(c, nil)
		return
	}

	// 返回favicon信息
	// 注意：缓存的图标 udId 为 -1，用于区分推荐图标
	response.Success(c, gin.H{
		"udId":        -1, // 缓存的图标使用 -1 表示，与推荐图标的正整数ID区分
		"title":       favicon.Title,
		"description": favicon.Description,
		"url":         url,
		"imgUrl":      favicon.ImgURL,
		"bgColor":     favicon.BgColor,
		"mimeType":    favicon.MimeType,
	})
}
