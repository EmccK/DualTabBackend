package admin

import (
	"dualtab-backend/internal/model"
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	repo *repository.CategoryRepo
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(repo *repository.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

// List 获取分类列表
func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.repo.FindAll(false)
	if err != nil {
		response.InternalError(c, "获取分类列表失败")
		return
	}
	response.Success(c, gin.H{"list": categories})
}

// Get 获取单个分类
func (h *CategoryHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	category, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.NotFound(c, "分类不存在")
		return
	}

	response.Success(c, category)
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	NameEn    string `json:"name_en"`
	SortOrder int    `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
}

// Create 创建分类
func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入分类名称")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	category := &model.Category{
		Name:      req.Name,
		NameEn:    req.NameEn,
		SortOrder: req.SortOrder,
		IsActive:  isActive,
	}

	if err := h.repo.Create(category); err != nil {
		response.InternalError(c, "创建分类失败")
		return
	}

	response.Success(c, category)
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name      string `json:"name"`
	NameEn    string `json:"name_en"`
	SortOrder *int   `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
}

// Update 更新分类
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	category, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.NotFound(c, "分类不存在")
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.NameEn != "" {
		category.NameEn = req.NameEn
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := h.repo.Update(category); err != nil {
		response.InternalError(c, "更新分类失败")
		return
	}

	response.Success(c, category)
}

// Delete 删除分类
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除分类失败")
		return
	}

	response.Success(c, nil)
}
