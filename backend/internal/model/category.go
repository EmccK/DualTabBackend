package model

import (
	"time"
)

// Category 图标分类
type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:50;not null"`
	NameEn    string    `json:"name_en" gorm:"column:name_en;size:50"`
	SortOrder int       `json:"sort_order" gorm:"column:sort_order;default:0"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active;default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 多对多关联
	Icons []Icon `json:"icons,omitempty" gorm:"many2many:icon_categories;"`
}

func (Category) TableName() string {
	return "categories"
}
