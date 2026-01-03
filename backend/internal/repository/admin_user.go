package repository

import (
	"dualtab-backend/internal/model"

	"gorm.io/gorm"
)

// AdminUserRepo 管理员用户仓库
type AdminUserRepo struct {
	db *gorm.DB
}

// NewAdminUserRepo 创建管理员用户仓库
func NewAdminUserRepo(db *gorm.DB) *AdminUserRepo {
	return &AdminUserRepo{db: db}
}

// FindByUsername 根据用户名查找
func (r *AdminUserRepo) FindByUsername(username string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据 ID 查找
func (r *AdminUserRepo) FindByID(id uint) (*model.AdminUser, error) {
	var user model.AdminUser
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *AdminUserRepo) Create(user *model.AdminUser) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *AdminUserRepo) Update(user *model.AdminUser) error {
	return r.db.Save(user).Error
}
