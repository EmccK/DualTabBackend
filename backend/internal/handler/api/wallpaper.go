package api

import (
	"dualtab-backend/internal/model"
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/response"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// WallpaperHandler 壁纸 API 处理器
type WallpaperHandler struct {
	repo       *repository.WallpaperRepo
	configRepo *repository.SystemConfigRepo
	client     *http.Client
}

// NewWallpaperHandler 创建壁纸 API 处理器
func NewWallpaperHandler(repo *repository.WallpaperRepo, configRepo *repository.SystemConfigRepo) *WallpaperHandler {
	return &WallpaperHandler{
		repo:       repo,
		configRepo: configRepo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetRandom 获取随机壁纸
// GET /wallpaper/random
func (h *WallpaperHandler) GetRandom(c *gin.Context) {
	// 先尝试从数据库获取
	wallpaper, err := h.repo.FindRandom()
	if err == nil && wallpaper != nil {
		response.Success(c, model.WallpaperResponse{
			ID:       wallpaper.UUID,
			Title:    wallpaper.Title,
			URL:      wallpaper.URL,
			ThumbURL: wallpaper.ThumbURL,
		})
		return
	}

	// 如果数据库没有壁纸，检查是否启用 Bing 壁纸
	bingEnabled := h.configRepo.GetValue(model.ConfigKeyBingWallpaperOn)
	if bingEnabled == "true" || bingEnabled == "1" {
		bingWallpaper, err := h.getBingWallpaper()
		if err == nil {
			response.Success(c, bingWallpaper)
			return
		}
	}

	// 返回默认壁纸
	response.Success(c, model.WallpaperResponse{
		ID:    "default",
		Title: "Default Wallpaper",
		URL:   "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=1920",
	})
}

// GetList 获取壁纸列表
// GET /wallpaper/list
func (h *WallpaperHandler) GetList(c *gin.Context) {
	wallpapers, err := h.repo.FindAll(true)
	if err != nil {
		response.Success(c, gin.H{"list": []interface{}{}})
		return
	}

	list := make([]model.WallpaperResponse, len(wallpapers))
	for i, w := range wallpapers {
		list[i] = model.WallpaperResponse{
			ID:       w.UUID,
			Title:    w.Title,
			URL:      w.URL,
			ThumbURL: w.ThumbURL,
		}
	}

	response.Success(c, gin.H{"list": list})
}

// getBingWallpaper 获取 Bing 每日壁纸
func (h *WallpaperHandler) getBingWallpaper() (*model.WallpaperResponse, error) {
	resp, err := h.client.Get("https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Images []struct {
			URL       string `json:"url"`
			Title     string `json:"title"`
			Copyright string `json:"copyright"`
		} `json:"images"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Images) == 0 {
		return nil, err
	}

	img := result.Images[0]
	return &model.WallpaperResponse{
		ID:    "bing-daily",
		Title: img.Copyright,
		URL:   "https://www.bing.com" + img.URL,
	}, nil
}
