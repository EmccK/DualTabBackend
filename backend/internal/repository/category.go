package repository

import (
	"dualtab-backend/internal/model"

	"gorm.io/gorm"
)

// CategoryRepo 分类仓库
type CategoryRepo struct {
	db *gorm.DB
}

// NewCategoryRepo 创建分类仓库
func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

// FindAll 查找所有分类
func (r *CategoryRepo) FindAll(onlyActive bool) ([]model.Category, error) {
	var categories []model.Category
	query := r.db.Order("sort_order ASC")
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&categories).Error
	return categories, err
}

// FindByID 根据 ID 查找
func (r *CategoryRepo) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Create 创建分类
func (r *CategoryRepo) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// Update 更新分类
func (r *CategoryRepo) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepo) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}
