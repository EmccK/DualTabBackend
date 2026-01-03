package model

import (
	"time"
)

// User 扩展用户（用于数据同步）
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:100;not null"`
	Name         string    `json:"name" gorm:"size:50"`
	Avatar       string    `json:"avatar" gorm:"size:500"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;size:255;not null"`
	Secret       string    `json:"secret" gorm:"uniqueIndex;size:100"` // 用于数据同步的 token
	IsBackup     int       `json:"isBackup" gorm:"default:0"`          // 是否已备份
	IsActivate   int       `json:"isActivate" gorm:"default:1"`        // 是否激活
	LastVisitAt  int64     `json:"lastVisitAt"`                        // 最后访问时间戳
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}

// UserData 用户数据（用于同步）
type UserData struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Type      string    `json:"type" gorm:"size:20;not null;index"` // icons, common, background, searcher, sidebar, todos, standby
	Data      string    `json:"data" gorm:"type:text"`              // JSON 数据
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserData) TableName() string {
	return "user_data"
}

// UserResponse 用户响应（不含敏感信息）
type UserResponse struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Secret      string `json:"secret"`
	IsBackup    int    `json:"isBackup"`
	LastVisitAt int64  `json:"lastVisitAt"`
	IsActivate  int    `json:"isActivate"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		Email:       u.Email,
		Name:        u.Name,
		Avatar:      u.Avatar,
		Secret:      u.Secret,
		IsBackup:    u.IsBackup,
		LastVisitAt: u.LastVisitAt,
		IsActivate:  u.IsActivate,
		CreatedAt:   u.CreatedAt.Unix(),
		UpdatedAt:   u.UpdatedAt.Unix(),
	}
}
