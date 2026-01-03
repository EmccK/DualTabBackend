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

// WallpaperHandler 壁纸管理处理器
type WallpaperHandler struct {
	repo     *repository.WallpaperRepo
	uploader *upload.Uploader
}

// NewWallpaperHandler 创建壁纸管理处理器
func NewWallpaperHandler(repo *repository.WallpaperRepo, uploader *upload.Uploader) *WallpaperHandler {
	return &WallpaperHandler{
		repo:     repo,
		uploader: uploader,
	}
}

// List 获取壁纸列表
func (h *WallpaperHandler) List(c *gin.Context) {
	wallpapers, err := h.repo.FindAll(false)
	if err != nil {
		response.InternalError(c, "获取壁纸列表失败")
		return
	}
	response.Success(c, gin.H{"list": wallpapers})
}

// Get 获取壁纸详情
func (h *WallpaperHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	wallpaper, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.NotFound(c, "壁纸不存在")
		return
	}

	response.Success(c, wallpaper)
}

// CreateWallpaperRequest 创建壁纸请求
type CreateWallpaperRequest struct {
	Title     string `json:"title"`
	URL       string `json:"url" binding:"required"`
	ThumbURL  string `json:"thumb_url"`
	Source    string `json:"source"`
	SortOrder int    `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
}

// Create 创建壁纸
func (h *WallpaperHandler) Create(c *gin.Context) {
	var req CreateWallpaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请提供壁纸 URL")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	wallpaper := &model.Wallpaper{
		UUID:      uuid.New().String(),
		Title:     req.Title,
		URL:       req.URL,
		ThumbURL:  req.ThumbURL,
		Source:    req.Source,
		SortOrder: req.SortOrder,
		IsActive:  isActive,
	}

	if err := h.repo.Create(wallpaper); err != nil {
		response.InternalError(c, "创建壁纸失败")
		return
	}

	response.Success(c, wallpaper)
}

// UpdateWallpaperRequest 更新壁纸请求
type UpdateWallpaperRequest struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	ThumbURL  string `json:"thumb_url"`
	Source    string `json:"source"`
	SortOrder *int   `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
}

// Update 更新壁纸
func (h *WallpaperHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	wallpaper, err := h.repo.FindByID(uint(id))
	if err != nil {
		response.NotFound(c, "壁纸不存在")
		return
	}

	var req UpdateWallpaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.Title != "" {
		wallpaper.Title = req.Title
	}
	if req.URL != "" {
		wallpaper.URL = req.URL
	}
	if req.ThumbURL != "" {
		wallpaper.ThumbURL = req.ThumbURL
	}
	if req.Source != "" {
		wallpaper.Source = req.Source
	}
	if req.SortOrder != nil {
		wallpaper.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		wallpaper.IsActive = *req.IsActive
	}

	if err := h.repo.Update(wallpaper); err != nil {
		response.InternalError(c, "更新壁纸失败")
		return
	}

	response.Success(c, wallpaper)
}

// Delete 删除壁纸
func (h *WallpaperHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除壁纸失败")
		return
	}

	response.Success(c, nil)
}

// UploadWallpaper 上传壁纸
func (h *WallpaperHandler) UploadWallpaper(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 上传文件
	fileURL, err := h.uploader.Upload(file, "wallpapers")
	if err != nil {
		response.InternalError(c, "上传失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"url": fileURL})
}
