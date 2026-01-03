package admin

import (
	"dualtab-backend/internal/model"
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchEngineHandler 搜索引擎处理器
type SearchEngineHandler struct {
	repo *repository.SearchEngineRepo
}

// NewSearchEngineHandler 创建搜索引擎处理器
func NewSearchEngineHandler(repo *repository.SearchEngineRepo) *SearchEngineHandler {
	return &SearchEngineHandler{repo: repo}
}

// List 获取搜索引擎列表
func (h *SearchEngineHandler) List(c *gin.Context) {
	engines, err := h.repo.FindAll(false)
	if err != nil {
		response.InternalError(c, "获取搜索引擎列表失败")
		return
	}
	response.Success(c, gin.H{"list": engines})
}

// Get 获取单个搜索引擎
func (h *SearchEngineHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	engine, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.NotFound(c, "搜索引擎不存在")
		return
	}

	response.Success(c, engine)
}

// CreateSearchEngineRequest 创建搜索引擎请求
type CreateSearchEngineRequest struct {
	Name      string `json:"name" binding:"required"`
	URL       string `json:"url" binding:"required"`
	IconURL   string `json:"icon_url"`
	SortOrder int    `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
}

// Create 创建搜索引擎
func (h *SearchEngineHandler) Create(c *gin.Context) {
	var req CreateSearchEngineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入名称和 URL")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	engine := &model.SearchEngine{
		UUID:      uuid.New().String(),
		Name:      req.Name,
		URL:       req.URL,
		IconURL:   req.IconURL,
		SortOrder: req.SortOrder,
		IsActive:  isActive,
	}

	if err := h.repo.Create(engine); err != nil {
		response.InternalError(c, "创建搜索引擎失败")
		return
	}

	response.Success(c, engine)
}

// UpdateSearchEngineRequest 更新搜索引擎请求
type UpdateSearchEngineRequest struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	IconURL   string `json:"icon_url"`
	SortOrder *int   `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
}

// Update 更新搜索引擎
func (h *SearchEngineHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	engine, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.NotFound(c, "搜索引擎不存在")
		return
	}

	var req UpdateSearchEngineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.Name != "" {
		engine.Name = req.Name
	}
	if req.URL != "" {
		engine.URL = req.URL
	}
	if req.IconURL != "" {
		engine.IconURL = req.IconURL
	}
	if req.SortOrder != nil {
		engine.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		engine.IsActive = *req.IsActive
	}

	if err := h.repo.Update(engine); err != nil {
		response.InternalError(c, "更新搜索引擎失败")
		return
	}

	response.Success(c, engine)
}

// Delete 删除搜索引擎
func (h *SearchEngineHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除搜索引擎失败")
		return
	}

	response.Success(c, nil)
}
