package model

import (
	"time"
)

// SearchEngine 搜索引擎
type SearchEngine struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UUID      string    `json:"uuid" gorm:"uniqueIndex;size:36;not null"`
	Name      string    `json:"name" gorm:"size:50;not null"`
	URL       string    `json:"url" gorm:"size:500;not null"` // 搜索 URL，%s 为占位符
	IconURL   string    `json:"icon_url" gorm:"column:icon_url;size:500"`
	SortOrder int       `json:"sort_order" gorm:"column:sort_order;default:0"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active;default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SearchEngine) TableName() string {
	return "search_engines"
}

// SearchEngineResponse 对外 API 响应格式
type SearchEngineResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
}

// ToResponse 转换为对外响应格式
func (s *SearchEngine) ToResponse() SearchEngineResponse {
	return SearchEngineResponse{
		ID:   s.UUID,
		Name: s.Name,
		URL:  s.URL,
		Icon: s.IconURL,
	}
}
