package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// WeatherService 天气服务
type WeatherService struct {
	apiKey  string
	apiType string // qweather 或 openweather
	client  *http.Client
}

// NewWeatherService 创建天气服务
func NewWeatherService(apiKey, apiType string) *WeatherService {
	return &WeatherService{
		apiKey:  apiKey,
		apiType: apiType,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// LocationResult 位置搜索结果
type LocationResult struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Adm1    string `json:"adm1"` // 省/州
	Adm2    string `json:"adm2"` // 市
	Lat     string `json:"lat"`
	Lon     string `json:"lon"`
}

// WeatherResult 天气结果
type WeatherResult struct {
	Location       string `json:"location"`
	Temperature    string `json:"temperature"`    // 当前温度
	TempMax        string `json:"tempMax"`        // 最高温度
	TempMin        string `json:"tempMin"`        // 最低温度
	FeelsLike      string `json:"feelsLike"`
	Text           string `json:"text"`        // 天气描述
	Icon           string `json:"icon"`        // 天气图标代码
	Humidity       string `json:"humidity"`    // 湿度
	WindDir        string `json:"windDir"`     // 风向
	WindScale      string `json:"windScale"`   // 风力等级
	UpdateTime     string `json:"updateTime"`
}

// SearchLocations 搜索城市位置
func (s *WeatherService) SearchLocations(keyword string) ([]LocationResult, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("天气 API Key 未配置")
	}

	if s.apiType == "openweather" {
		return s.searchLocationsOpenWeather(keyword)
	}
	// 默认使用和风天气
	return s.searchLocationsQWeather(keyword)
}

// GetWeather 获取天气信息
func (s *WeatherService) GetWeather(locationID string) (*WeatherResult, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("天气 API Key 未配置")
	}

	if s.apiType == "openweather" {
		return s.getWeatherOpenWeather(locationID)
	}
	// 默认使用和风天气
	return s.getWeatherQWeather(locationID)
}

// ========== 和风天气 API ==========

func (s *WeatherService) searchLocationsQWeather(keyword string) ([]LocationResult, error) {
	apiURL := fmt.Sprintf(
		"https://geoapi.qweather.com/v2/city/lookup?location=%s&key=%s",
		url.QueryEscape(keyword),
		s.apiKey,
	)

	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code     string `json:"code"`
		Location []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Country string `json:"country"`
			Adm1    string `json:"adm1"`
			Adm2    string `json:"adm2"`
			Lat     string `json:"lat"`
			Lon     string `json:"lon"`
		} `json:"location"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != "200" {
		return nil, fmt.Errorf("和风天气 API 错误: %s", result.Code)
	}

	locations := make([]LocationResult, len(result.Location))
	for i, loc := range result.Location {
		locations[i] = LocationResult{
			ID:      loc.ID,
			Name:    loc.Name,
			Country: loc.Country,
			Adm1:    loc.Adm1,
			Adm2:    loc.Adm2,
			Lat:     loc.Lat,
			Lon:     loc.Lon,
		}
	}

	return locations, nil
}

func (s *WeatherService) getWeatherQWeather(locationID string) (*WeatherResult, error) {
	// 1. 获取实时天气
	nowURL := fmt.Sprintf(
		"https://devapi.qweather.com/v7/weather/now?location=%s&key=%s",
		url.QueryEscape(locationID),
		s.apiKey,
	)

	nowResp, err := s.client.Get(nowURL)
	if err != nil {
		return nil, err
	}
	defer nowResp.Body.Close()

	nowBody, err := io.ReadAll(nowResp.Body)
	if err != nil {
		return nil, err
	}

	var nowResult struct {
		Code string `json:"code"`
		Now  struct {
			Temp      string `json:"temp"`
			FeelsLike string `json:"feelsLike"`
			Text      string `json:"text"`
			Icon      string `json:"icon"`
			Humidity  string `json:"humidity"`
			WindDir   string `json:"windDir"`
			WindScale string `json:"windScale"`
		} `json:"now"`
		UpdateTime string `json:"updateTime"`
	}

	if err := json.Unmarshal(nowBody, &nowResult); err != nil {
		return nil, err
	}

	if nowResult.Code != "200" {
		return nil, fmt.Errorf("和风天气 API 错误: %s", nowResult.Code)
	}

	// 2. 获取每日预报（获取今日最高最低温度）
	dailyURL := fmt.Sprintf(
		"https://devapi.qweather.com/v7/weather/3d?location=%s&key=%s",
		url.QueryEscape(locationID),
		s.apiKey,
	)

	dailyResp, err := s.client.Get(dailyURL)
	if err != nil {
		// 如果每日预报失败，仍返回实时天气，温度范围留空
		return &WeatherResult{
			Temperature: nowResult.Now.Temp,
			TempMax:     "",
			TempMin:     "",
			FeelsLike:   nowResult.Now.FeelsLike,
			Text:        nowResult.Now.Text,
			Icon:        nowResult.Now.Icon,
			Humidity:    nowResult.Now.Humidity,
			WindDir:     nowResult.Now.WindDir,
			WindScale:   nowResult.Now.WindScale,
			UpdateTime:  nowResult.UpdateTime,
		}, nil
	}
	defer dailyResp.Body.Close()

	dailyBody, err := io.ReadAll(dailyResp.Body)
	if err != nil {
		return &WeatherResult{
			Temperature: nowResult.Now.Temp,
			TempMax:     "",
			TempMin:     "",
			FeelsLike:   nowResult.Now.FeelsLike,
			Text:        nowResult.Now.Text,
			Icon:        nowResult.Now.Icon,
			Humidity:    nowResult.Now.Humidity,
			WindDir:     nowResult.Now.WindDir,
			WindScale:   nowResult.Now.WindScale,
			UpdateTime:  nowResult.UpdateTime,
		}, nil
	}

	var dailyResult struct {
		Code  string `json:"code"`
		Daily []struct {
			TempMax string `json:"tempMax"`
			TempMin string `json:"tempMin"`
		} `json:"daily"`
	}

	if err := json.Unmarshal(dailyBody, &dailyResult); err != nil || dailyResult.Code != "200" || len(dailyResult.Daily) == 0 {
		return &WeatherResult{
			Temperature: nowResult.Now.Temp,
			TempMax:     "",
			TempMin:     "",
			FeelsLike:   nowResult.Now.FeelsLike,
			Text:        nowResult.Now.Text,
			Icon:        nowResult.Now.Icon,
			Humidity:    nowResult.Now.Humidity,
			WindDir:     nowResult.Now.WindDir,
			WindScale:   nowResult.Now.WindScale,
			UpdateTime:  nowResult.UpdateTime,
		}, nil
	}

	// 返回包含今日最高最低温度的天气数据
	return &WeatherResult{
		Temperature: nowResult.Now.Temp,
		TempMax:     dailyResult.Daily[0].TempMax,
		TempMin:     dailyResult.Daily[0].TempMin,
		FeelsLike:   nowResult.Now.FeelsLike,
		Text:        nowResult.Now.Text,
		Icon:        nowResult.Now.Icon,
		Humidity:    nowResult.Now.Humidity,
		WindDir:     nowResult.Now.WindDir,
		WindScale:   nowResult.Now.WindScale,
		UpdateTime:  nowResult.UpdateTime,
	}, nil
}

// ========== OpenWeather API ==========

func (s *WeatherService) searchLocationsOpenWeather(keyword string) ([]LocationResult, error) {
	apiURL := fmt.Sprintf(
		"http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=10&appid=%s",
		url.QueryEscape(keyword),
		s.apiKey,
	)

	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []struct {
		Name    string  `json:"name"`
		Country string  `json:"country"`
		State   string  `json:"state"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	}

	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	locations := make([]LocationResult, len(results))
	for i, loc := range results {
		locations[i] = LocationResult{
			ID:      fmt.Sprintf("%.4f,%.4f", loc.Lat, loc.Lon),
			Name:    loc.Name,
			Country: loc.Country,
			Adm1:    loc.State,
			Lat:     fmt.Sprintf("%.4f", loc.Lat),
			Lon:     fmt.Sprintf("%.4f", loc.Lon),
		}
	}

	return locations, nil
}

func (s *WeatherService) getWeatherOpenWeather(locationID string) (*WeatherResult, error) {
	// locationID 格式: "lat,lon"
	var lat, lon string
	fmt.Sscanf(locationID, "%[^,],%s", &lat, &lon)

	// 1. 获取当前天气
	currentURL := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=metric&lang=zh_cn&appid=%s",
		lat, lon, s.apiKey,
	)

	currentResp, err := s.client.Get(currentURL)
	if err != nil {
		return nil, err
	}
	defer currentResp.Body.Close()

	currentBody, err := io.ReadAll(currentResp.Body)
	if err != nil {
		return nil, err
	}

	var currentResult struct {
		Weather []struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Humidity  int     `json:"humidity"`
		} `json:"main"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		} `json:"wind"`
		Name string `json:"name"`
		Dt   int64  `json:"dt"`
	}

	if err := json.Unmarshal(currentBody, &currentResult); err != nil {
		return nil, err
	}

	// 风向转换
	windDir := getWindDirection(currentResult.Wind.Deg)
	// 风力等级转换 (m/s to scale)
	windScale := getWindScale(currentResult.Wind.Speed)

	weather := &WeatherResult{
		Location:    currentResult.Name,
		Temperature: fmt.Sprintf("%.0f", currentResult.Main.Temp),
		TempMax:     fmt.Sprintf("%.0f", currentResult.Main.TempMax),
		TempMin:     fmt.Sprintf("%.0f", currentResult.Main.TempMin),
		FeelsLike:   fmt.Sprintf("%.0f", currentResult.Main.FeelsLike),
		Humidity:    fmt.Sprintf("%d", currentResult.Main.Humidity),
		WindDir:     windDir,
		WindScale:   windScale,
		UpdateTime:  time.Unix(currentResult.Dt, 0).Format(time.RFC3339),
	}

	if len(currentResult.Weather) > 0 {
		weather.Text = currentResult.Weather[0].Description
		weather.Icon = currentResult.Weather[0].Icon
	}

	return weather, nil
}

// getWindDirection 根据角度获取风向
func getWindDirection(deg int) string {
	directions := []string{"北", "东北", "东", "东南", "南", "西南", "西", "西北"}
	index := ((deg + 22) % 360) / 45
	return directions[index]
}

// getWindScale 根据风速获取风力等级
func getWindScale(speed float64) string {
	switch {
	case speed < 0.3:
		return "0"
	case speed < 1.6:
		return "1"
	case speed < 3.4:
		return "2"
	case speed < 5.5:
		return "3"
	case speed < 8.0:
		return "4"
	case speed < 10.8:
		return "5"
	case speed < 13.9:
		return "6"
	case speed < 17.2:
		return "7"
	case speed < 20.8:
		return "8"
	case speed < 24.5:
		return "9"
	case speed < 28.5:
		return "10"
	case speed < 32.7:
		return "11"
	default:
		return "12"
	}
}
