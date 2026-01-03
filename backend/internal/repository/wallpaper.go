package repository

import (
	"dualtab-backend/internal/model"

	"gorm.io/gorm"
)

// WallpaperRepo 壁纸仓库
type WallpaperRepo struct {
	db *gorm.DB
}

// NewWallpaperRepo 创建壁纸仓库
func NewWallpaperRepo(db *gorm.DB) *WallpaperRepo {
	return &WallpaperRepo{db: db}
}

// FindAll 获取所有壁纸
func (r *WallpaperRepo) FindAll(activeOnly bool) ([]model.Wallpaper, error) {
	var wallpapers []model.Wallpaper
	query := r.db.Order("sort_order ASC, id DESC")
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&wallpapers).Error
	return wallpapers, err
}

// FindByID 根据 ID 获取壁纸
func (r *WallpaperRepo) FindByID(id uint) (*model.Wallpaper, error) {
	var wallpaper model.Wallpaper
	err := r.db.First(&wallpaper, id).Error
	return &wallpaper, err
}

// FindRandom 随机获取一张壁纸
func (r *WallpaperRepo) FindRandom() (*model.Wallpaper, error) {
	var wallpaper model.Wallpaper
	err := r.db.Where("is_active = ?", true).Order("RANDOM()").First(&wallpaper).Error
	return &wallpaper, err
}

// Create 创建壁纸
func (r *WallpaperRepo) Create(wallpaper *model.Wallpaper) error {
	return r.db.Create(wallpaper).Error
}

// Update 更新壁纸
func (r *WallpaperRepo) Update(wallpaper *model.Wallpaper) error {
	return r.db.Save(wallpaper).Error
}

// Delete 删除壁纸
func (r *WallpaperRepo) Delete(id uint) error {
	return r.db.Delete(&model.Wallpaper{}, id).Error
}

// Count 统计壁纸数量
func (r *WallpaperRepo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Wallpaper{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}
