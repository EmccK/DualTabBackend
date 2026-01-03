package repository

import (
	"dualtab-backend/internal/model"

	"gorm.io/gorm"
)

// SystemConfigRepo 系统配置仓库
type SystemConfigRepo struct {
	db *gorm.DB
}

// NewSystemConfigRepo 创建系统配置仓库
func NewSystemConfigRepo(db *gorm.DB) *SystemConfigRepo {
	return &SystemConfigRepo{db: db}
}

// FindAll 获取所有配置
func (r *SystemConfigRepo) FindAll() ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	err := r.db.Order("key ASC").Find(&configs).Error
	return configs, err
}

// FindByKey 根据 Key 获取配置
func (r *SystemConfigRepo) FindByKey(key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := r.db.Where("\"key\" = ?", key).First(&config).Error
	return &config, err
}

// GetValue 获取配置值
func (r *SystemConfigRepo) GetValue(key string) string {
	config, err := r.FindByKey(key)
	if err != nil {
		return ""
	}
	return config.Value
}

// SetValue 设置配置值
func (r *SystemConfigRepo) SetValue(key, value, remark string) error {
	var config model.SystemConfig
	err := r.db.Where("\"key\" = ?", key).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		// 创建新配置
		config = model.SystemConfig{
			Key:    key,
			Value:  value,
			Remark: remark,
		}
		return r.db.Create(&config).Error
	}
	if err != nil {
		return err
	}
	// 更新现有配置
	config.Value = value
	if remark != "" {
		config.Remark = remark
	}
	return r.db.Save(&config).Error
}

// Delete 删除配置
func (r *SystemConfigRepo) Delete(key string) error {
	return r.db.Where("\"key\" = ?", key).Delete(&model.SystemConfig{}).Error
}

// BatchSetValue 批量设置配置值（使用事务）
func (r *SystemConfigRepo) BatchSetValue(configs []struct {
	Key    string
	Value  string
	Remark string
}) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, cfg := range configs {
			var config model.SystemConfig
			err := tx.Where("\"key\" = ?", cfg.Key).First(&config).Error
			if err == gorm.ErrRecordNotFound {
				// 创建新配置
				config = model.SystemConfig{
					Key:    cfg.Key,
					Value:  cfg.Value,
					Remark: cfg.Remark,
				}
				if err := tx.Create(&config).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				// 更新现有配置
				config.Value = cfg.Value
				if cfg.Remark != "" {
					config.Remark = cfg.Remark
				}
				if err := tx.Save(&config).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
