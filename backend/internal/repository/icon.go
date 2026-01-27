package repository

import (
	"dualtab-backend/internal/model"

	"gorm.io/gorm"
)

// IconRepo 图标仓库
type IconRepo struct {
	db *gorm.DB
}

// NewIconRepo 创建图标仓库
func NewIconRepo(db *gorm.DB) *IconRepo {
	return &IconRepo{db: db}
}

// IconQuery 图标查询参数
type IconQuery struct {
	CategoryID uint
	Keyword    string
	IsActive   *bool
	Page       int
	Size       int
}

// FindAll 查找图标列表
func (r *IconRepo) FindAll(query IconQuery) ([]model.Icon, int64, error) {
	var icons []model.Icon
	var total int64

	db := r.db.Model(&model.Icon{})

	// 分类筛选：通过中间表 JOIN 查询
	if query.CategoryID > 0 {
		db = db.Joins("JOIN icon_categories ON icon_categories.icon_id = icons.id").
			Where("icon_categories.category_id = ?", query.CategoryID)
	}

	// 关键词搜索
	if query.Keyword != "" {
		db = db.Where("icons.title ILIKE ? OR icons.description ILIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 状态筛选
	if query.IsActive != nil {
		db = db.Where("icons.is_active = ?", *query.IsActive)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if query.Page > 0 && query.Size > 0 {
		offset := (query.Page - 1) * query.Size
		db = db.Offset(offset).Limit(query.Size)
	} else if query.Size > 0 {
		db = db.Limit(query.Size)
	}

	// 排序并查询，Preload Categories
	err := db.Order("icons.sort_order ASC, icons.id DESC").Preload("Categories").Find(&icons).Error
	return icons, total, err
}

// FindByID 根据 ID 查找
func (r *IconRepo) FindByID(id uint) (*model.Icon, error) {
	var icon model.Icon
	err := r.db.Preload("Categories").First(&icon, id).Error
	if err != nil {
		return nil, err
	}
	return &icon, nil
}

// FindByUUID 根据 UUID 查找
func (r *IconRepo) FindByUUID(uuid string) (*model.Icon, error) {
	var icon model.Icon
	err := r.db.Where("uuid = ?", uuid).First(&icon).Error
	if err != nil {
		return nil, err
	}
	return &icon, nil
}

// FindByURL 根据 URL 查找
func (r *IconRepo) FindByURL(url string) (*model.Icon, error) {
	var icon model.Icon
	err := r.db.Where("url = ? AND is_active = ?", url, true).First(&icon).Error
	if err != nil {
		return nil, err
	}
	return &icon, nil
}

// Create 创建图标
func (r *IconRepo) Create(icon *model.Icon) error {
	return r.db.Create(icon).Error
}

// Update 更新图标
func (r *IconRepo) Update(icon *model.Icon) error {
	return r.db.Save(icon).Error
}

// Delete 删除图标
func (r *IconRepo) Delete(id uint) error {
	return r.db.Delete(&model.Icon{}, id).Error
}

// UpdateCategories 更新图标的分类关联
func (r *IconRepo) UpdateCategories(icon *model.Icon, categoryIDs []uint) error {
	// 构建分类列表
	var categories []model.Category
	if len(categoryIDs) > 0 {
		if err := r.db.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
			return err
		}
	}
	// 替换关联
	return r.db.Model(icon).Association("Categories").Replace(categories)
}
