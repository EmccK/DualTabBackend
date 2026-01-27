package model

import (
	"time"
)

// IconCategory 图标-分类多对多关联表
type IconCategory struct {
	IconID     uint      `json:"icon_id" gorm:"primaryKey"`
	CategoryID uint      `json:"category_id" gorm:"primaryKey"`
	CreatedAt  time.Time `json:"created_at"`
}

func (IconCategory) TableName() string {
	return "icon_categories"
}
