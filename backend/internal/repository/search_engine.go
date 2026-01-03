package repository

import (
	"dualtab-backend/internal/model"

	"gorm.io/gorm"
)

// SearchEngineRepo 搜索引擎仓库
type SearchEngineRepo struct {
	db *gorm.DB
}

// NewSearchEngineRepo 创建搜索引擎仓库
func NewSearchEngineRepo(db *gorm.DB) *SearchEngineRepo {
	return &SearchEngineRepo{db: db}
}

// FindAll 查找所有搜索引擎
func (r *SearchEngineRepo) FindAll(onlyActive bool) ([]model.SearchEngine, error) {
	var engines []model.SearchEngine
	query := r.db.Order("sort_order ASC")
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&engines).Error
	return engines, err
}

// FindByID 根据 ID 查找
func (r *SearchEngineRepo) FindByID(id uint) (*model.SearchEngine, error) {
	var engine model.SearchEngine
	err := r.db.First(&engine, id).Error
	if err != nil {
		return nil, err
	}
	return &engine, nil
}

// FindByUUID 根据 UUID 查找
func (r *SearchEngineRepo) FindByUUID(uuid string) (*model.SearchEngine, error) {
	var engine model.SearchEngine
	err := r.db.Where("uuid = ?", uuid).First(&engine).Error
	if err != nil {
		return nil, err
	}
	return &engine, nil
}

// Create 创建搜索引擎
func (r *SearchEngineRepo) Create(engine *model.SearchEngine) error {
	return r.db.Create(engine).Error
}

// Update 更新搜索引擎
func (r *SearchEngineRepo) Update(engine *model.SearchEngine) error {
	return r.db.Save(engine).Error
}

// Delete 删除搜索引擎
func (r *SearchEngineRepo) Delete(id uint) error {
	return r.db.Delete(&model.SearchEngine{}, id).Error
}
