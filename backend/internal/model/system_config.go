package model

import (
	"time"
)

// SystemConfig 系统配置
type SystemConfig struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"uniqueIndex;size:50;not null"`
	Value     string    `json:"value" gorm:"type:text"`
	Remark    string    `json:"remark" gorm:"size:200"` // 配置说明
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}

// 预定义配置 Key
const (
	ConfigKeyWeatherAPIKey   = "weather_api_key"   // 天气 API Key (和风天气)
	ConfigKeyWeatherAPIType  = "weather_api_type"  // 天气 API 类型: qweather/openweather
	ConfigKeySearchSuggestOn = "search_suggest_on" // 是否启用搜索建议代理
	ConfigKeyBingWallpaperOn = "bing_wallpaper_on" // 是否启用 Bing 每日壁纸
)
