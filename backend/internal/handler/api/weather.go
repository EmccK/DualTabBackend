package api

import (
	"dualtab-backend/internal/repository"
	"dualtab-backend/internal/service"
	"dualtab-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// WeatherHandler 天气 API 处理器
type WeatherHandler struct {
	configRepo *repository.SystemConfigRepo
}

// NewWeatherHandler 创建天气 API 处理器
func NewWeatherHandler(configRepo *repository.SystemConfigRepo) *WeatherHandler {
	return &WeatherHandler{configRepo: configRepo}
}

// getWeatherService 获取天气服务实例
func (h *WeatherHandler) getWeatherService() *service.WeatherService {
	apiKey := h.configRepo.GetValue("weather_api_key")
	apiType := h.configRepo.GetValue("weather_api_type")
	if apiType == "" {
		apiType = "qweather"
	}
	return service.NewWeatherService(apiKey, apiType)
}

// SearchLocations 搜索城市位置
// GET /weather/locations?keyword=北京
func (h *WeatherHandler) SearchLocations(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "请输入搜索关键词")
		return
	}

	weatherSvc := h.getWeatherService()
	locations, err := weatherSvc.SearchLocations(keyword)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"list": locations})
}

// GetWeather 获取天气信息
// GET /weather?location=101010100
func (h *WeatherHandler) GetWeather(c *gin.Context) {
	location := c.Query("location")
	if location == "" {
		response.BadRequest(c, "请提供位置 ID")
		return
	}

	weatherSvc := h.getWeatherService()
	weather, err := weatherSvc.GetWeather(location)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, weather)
}
