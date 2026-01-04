package repository

import (
	"dualtab-backend/internal/model"

	"gorm.io/gorm"
)

// FaviconCacheRepo 网站图标缓存仓库
type FaviconCacheRepo struct {
	db *gorm.DB
}

// NewFaviconCacheRepo 创建网站图标缓存仓库
func NewFaviconCacheRepo(db *gorm.DB) *FaviconCacheRepo {
	return &FaviconCacheRepo{db: db}
}

// FindByHost 根据域名查找缓存的图标
func (r *FaviconCacheRepo) FindByHost(host string) (*model.FaviconCache, error) {
	var cache model.FaviconCache
	err := r.db.Where("host = ?", host).First(&cache).Error
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// Create 创建缓存记录
func (r *FaviconCacheRepo) Create(cache *model.FaviconCache) error {
	return r.db.Create(cache).Error
}

// Update 更新缓存记录
func (r *FaviconCacheRepo) Update(cache *model.FaviconCache) error {
	return r.db.Save(cache).Error
}
