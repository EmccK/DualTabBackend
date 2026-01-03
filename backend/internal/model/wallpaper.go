package model

import (
	"time"
)

// Wallpaper 壁纸
type Wallpaper struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UUID      string    `json:"uuid" gorm:"uniqueIndex;size:36;not null"`
	Title     string    `json:"title" gorm:"size:100"`
	URL       string    `json:"url" gorm:"size:500;not null"`       // 壁纸图片 URL
	ThumbURL  string    `json:"thumb_url" gorm:"size:500"`          // 缩略图 URL
	Source    string    `json:"source" gorm:"size:50"`              // 来源：upload/unsplash/bing
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	IsActive  bool      `json:"is_active" gorm:"default:true;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Wallpaper) TableName() string {
	return "wallpapers"
}

// WallpaperResponse 对外 API 响应格式
type WallpaperResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	ThumbURL string `json:"thumbUrl"`
}
