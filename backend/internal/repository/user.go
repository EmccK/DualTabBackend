package repository

import (
	"dualtab-backend/internal/model"

	"gorm.io/gorm"
)

// UserRepo 用户仓库
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建用户仓库
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create 创建用户
func (r *UserRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindBySecret 根据 secret 查找用户
func (r *UserRepo) FindBySecret(secret string) (*model.User, error) {
	var user model.User
	err := r.db.Where("secret = ?", secret).First(&user).Error
	return &user, err
}

// FindByID 根据 ID 查找用户
func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// Update 更新用户
func (r *UserRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// UpdateLastVisit 更新最后访问时间
func (r *UserRepo) UpdateLastVisit(id uint, timestamp int64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("last_visit_at", timestamp).Error
}

// UserDataRepo 用户数据仓库
type UserDataRepo struct {
	db *gorm.DB
}

// NewUserDataRepo 创建用户数据仓库
func NewUserDataRepo(db *gorm.DB) *UserDataRepo {
	return &UserDataRepo{db: db}
}

// GetByType 获取指定类型的用户数据
func (r *UserDataRepo) GetByType(userID uint, dataType string) (*model.UserData, error) {
	var data model.UserData
	err := r.db.Where("user_id = ? AND type = ?", userID, dataType).First(&data).Error
	return &data, err
}

// GetAllByUserID 获取用户所有数据
func (r *UserDataRepo) GetAllByUserID(userID uint) ([]model.UserData, error) {
	var dataList []model.UserData
	err := r.db.Where("user_id = ?", userID).Find(&dataList).Error
	return dataList, err
}

// Upsert 创建或更新用户数据
func (r *UserDataRepo) Upsert(userID uint, dataType string, data string) error {
	var existing model.UserData
	err := r.db.Where("user_id = ? AND type = ?", userID, dataType).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		newData := model.UserData{
			UserID: userID,
			Type:   dataType,
			Data:   data,
		}
		return r.db.Create(&newData).Error
	}

	if err != nil {
		return err
	}

	// 更新现有记录
	existing.Data = data
	return r.db.Save(&existing).Error
}

// Delete 删除用户数据
func (r *UserDataRepo) Delete(userID uint, dataType string) error {
	return r.db.Where("user_id = ? AND type = ?", userID, dataType).Delete(&model.UserData{}).Error
}
