package model

import (
	"time"
)

// Icon 推荐图标
type Icon struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UUID        string    `json:"uuid" gorm:"uniqueIndex;size:36;not null"`
	Title       string    `json:"title" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	URL         string    `json:"url" gorm:"size:500;not null"`
	ImgURL      string    `json:"img_url" gorm:"column:img_url;size:500"`
	BgColor     string    `json:"bg_color" gorm:"column:bg_color;size:20;default:'#ffffff'"`
	MimeType    string    `json:"mime_type" gorm:"column:mime_type;size:50;default:'image/png'"`
	SortOrder   int       `json:"sort_order" gorm:"column:sort_order;default:0"`
	IsActive    bool      `json:"is_active" gorm:"column:is_active;default:true;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 多对多关联
	Categories []Category `json:"categories,omitempty" gorm:"many2many:icon_categories;"`
}

func (Icon) TableName() string {
	return "icons"
}

// IconResponse 对外 API 响应格式（兼容 MonkNow）
type IconResponse struct {
	UdID        uint   `json:"udId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ImgURL      string `json:"imgUrl"`
	BgColor     string `json:"bgColor"`
	MimeType    string `json:"mimeType"`
}

// ToResponse 转换为对外响应格式
func (i *Icon) ToResponse() IconResponse {
	return IconResponse{
		UdID:        i.ID,
		Title:       i.Title,
		Description: i.Description,
		URL:         i.URL,
		ImgURL:      i.ImgURL,
		BgColor:     i.BgColor,
		MimeType:    i.MimeType,
	}
}
