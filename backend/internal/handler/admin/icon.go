package admin

import (
	"dualtab-backend/internal/model"
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/response"
	"dualtab-backend/pkg/upload"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IconHandler 图标处理器
type IconHandler struct {
	repo     *repository.IconRepo
	uploader *upload.Uploader
}

// NewIconHandler 创建图标处理器
func NewIconHandler(repo *repository.IconRepo, uploader *upload.Uploader) *IconHandler {
	return &IconHandler{
		repo:     repo,
		uploader: uploader,
	}
}

// List 获取图标列表
func (h *IconHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 32)
	keyword := c.Query("keyword")

	query := repository.IconQuery{
		CategoryID: uint(categoryID),
		Keyword:    keyword,
		Page:       page,
		Size:       size,
	}

	icons, total, err := h.repo.FindAll(query)
	if err != nil {
		response.InternalError(c, "获取图标列表失败")
		return
	}

	response.SuccessWithPage(c, icons, total, page, size)
}

// Get 获取单个图标
func (h *IconHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	icon, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.NotFound(c, "图标不存在")
		return
	}

	response.Success(c, icon)
}

// CreateIconRequest 创建图标请求
type CreateIconRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"required"`
	ImgURL      string `json:"img_url"`
	BgColor     string `json:"bg_color"`
	MimeType    string `json:"mime_type"`
	CategoryID  uint   `json:"category_id"`
	SortOrder   int    `json:"sort_order"`
	IsActive    *bool  `json:"is_active"`
}

// Create 创建图标
func (h *IconHandler) Create(c *gin.Context) {
	var req CreateIconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入标题和 URL")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	bgColor := "#ffffff"
	if req.BgColor != "" {
		bgColor = req.BgColor
	}

	mimeType := "image/png"
	if req.MimeType != "" {
		mimeType = req.MimeType
	}

	icon := &model.Icon{
		UUID:        uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		ImgURL:      req.ImgURL,
		BgColor:     bgColor,
		MimeType:    mimeType,
		CategoryID:  req.CategoryID,
		SortOrder:   req.SortOrder,
		IsActive:    isActive,
	}

	if err := h.repo.Create(icon); err != nil {
		response.InternalError(c, "创建图标失败")
		return
	}

	response.Success(c, icon)
}

// UpdateIconRequest 更新图标请求
type UpdateIconRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ImgURL      string `json:"img_url"`
	BgColor     string `json:"bg_color"`
	MimeType    string `json:"mime_type"`
	CategoryID  *uint  `json:"category_id"`
	SortOrder   *int   `json:"sort_order"`
	IsActive    *bool  `json:"is_active"`
}

// Update 更新图标
func (h *IconHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	icon, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.NotFound(c, "图标不存在")
		return
	}

	var req UpdateIconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.Title != "" {
		icon.Title = req.Title
	}
	if req.Description != "" {
		icon.Description = req.Description
	}
	if req.URL != "" {
		icon.URL = req.URL
	}
	if req.ImgURL != "" {
		icon.ImgURL = req.ImgURL
	}
	if req.BgColor != "" {
		icon.BgColor = req.BgColor
	}
	if req.MimeType != "" {
		icon.MimeType = req.MimeType
	}
	if req.CategoryID != nil {
		icon.CategoryID = *req.CategoryID
	}
	if req.SortOrder != nil {
		icon.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		icon.IsActive = *req.IsActive
	}

	if err := h.repo.Update(icon); err != nil {
		response.InternalError(c, "更新图标失败")
		return
	}

	response.Success(c, icon)
}

// Delete 删除图标
func (h *IconHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除图标失败")
		return
	}

	response.Success(c, nil)
}

// UploadIcon 上传图标图片
func (h *IconHandler) UploadIcon(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择文件")
		return
	}

	url, err := h.uploader.UploadIcon(file)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{"url": url})
}
