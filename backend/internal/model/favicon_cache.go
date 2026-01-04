package model

import (
	"time"
)

// FaviconCache 网站图标缓存表
// 用于缓存从外部API获取的网站图标，避免重复请求
type FaviconCache struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Host        string    `json:"host" gorm:"uniqueIndex;size:255;not null;comment:网站域名(如: example.com)"`
	Title       string    `json:"title" gorm:"size:200;comment:网站标题"`
	Description string    `json:"description" gorm:"type:text;comment:网站描述"`
	ImgURL      string    `json:"img_url" gorm:"column:img_url;size:500;comment:图标URL"`
	BgColor     string    `json:"bg_color" gorm:"column:bg_color;size:20;default:'#ffffff';comment:背景颜色"`
	MimeType    string    `json:"mime_type" gorm:"column:mime_type;size:50;default:'image/png';comment:图片类型"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (FaviconCache) TableName() string {
	return "favicon_cache"
}
