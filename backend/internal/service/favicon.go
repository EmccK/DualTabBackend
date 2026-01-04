package service

import (
	"dualtab-backend/internal/model"
	"dualtab-backend/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

// FaviconService 网站图标服务
type FaviconService struct {
	cacheRepo  *repository.FaviconCacheRepo
	httpClient *http.Client
	limiter    *rate.Limiter // API 限流器
	apiURL     string        // Favicon API URL
}

// NewFaviconService 创建网站图标服务
func NewFaviconService(cacheRepo *repository.FaviconCacheRepo, apiURL string) *FaviconService {
	return &FaviconService{
		cacheRepo: cacheRepo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				DisableKeepAlives:   false,
				DisableCompression:  false,
				MaxIdleConnsPerHost: 10,
			},
			// 限制重定向次数
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("重定向次数过多")
				}
				return nil
			},
		},
		// 限流器：每秒最多 10 个请求，突发最多 20 个
		limiter: rate.NewLimiter(rate.Every(100*time.Millisecond), 20),
		apiURL:  apiURL,
	}
}

// FaviconInfo 图标信息
type FaviconInfo struct {
	Title       string
	Description string
	ImgURL      string
	BgColor     string
	MimeType    string
}

// MonkNowAPIResponse MonkNow API 响应结构（简化版）
type MonkNowAPIResponse struct {
	Data struct {
		Icon struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			ImgURL      string `json:"imgUrl"`
			BgColor     string `json:"bgColor"`
			MimeType    string `json:"mimeType"`
		} `json:"icon"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// GetFavicon 获取网站图标
// 优先从数据库缓存获取
// 支持域名降级：例如 wiki.example.com -> example.com
func (s *FaviconService) GetFavicon(inputURL string) (*FaviconInfo, error) {
	// 解析URL获取host
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return nil, fmt.Errorf("无效的URL: %v", err)
	}

	// 验证协议，只允许 http 和 https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("不支持的协议: %s，仅支持 http 和 https", parsedURL.Scheme)
	}

	host := parsedURL.Host
	if host == "" {
		return nil, fmt.Errorf("无法从URL提取域名")
	}

	// 验证域名长度，防止过长的host导致数据库错误
	if len(host) > 255 {
		return nil, fmt.Errorf("域名长度超过限制")
	}

	// 尝试从缓存获取（支持域名降级）
	favicon, err := s.getFaviconFromCache(host)
	if err == nil && favicon != nil {
		log.Debug().Str("host", host).Msg("从缓存获取图标成功")
		return favicon, nil
	}

	// 缓存未命中，从外部API获取
	log.Info().Str("host", host).Msg("缓存未命中，尝试从外部API获取")
	favicon, err = s.fetchFaviconFromAPI(host)
	if err != nil {
		log.Warn().Err(err).Str("host", host).Msg("从外部API获取图标失败")
		return nil, fmt.Errorf("获取图标失败: %v", err)
	}

	// 保存到缓存
	s.saveFaviconToCache(host, favicon)

	return favicon, nil
}

// getFaviconFromCache 从缓存获取图标，支持域名降级
// 降级策略：只降级一次，从子域名降到父域名
// 例如: wiki.example.com -> example.com
// 不会继续降级到顶级域名(如 com)，因为顶级域名缓存没有实际意义
func (s *FaviconService) getFaviconFromCache(host string) (*FaviconInfo, error) {
	// 尝试当前域名
	cache, err := s.cacheRepo.FindByHost(host)
	if err == nil {
		return &FaviconInfo{
			Title:       cache.Title,
			Description: cache.Description,
			ImgURL:      cache.ImgURL,
			BgColor:     cache.BgColor,
			MimeType:    cache.MimeType,
		}, nil
	}

	// 尝试父域名（仅降级一次）
	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		// 去掉子域名，尝试父域名
		// wiki.example.com -> example.com
		parentHost := strings.Join(parts[1:], ".")
		cache, err = s.cacheRepo.FindByHost(parentHost)
		if err == nil {
			return &FaviconInfo{
				Title:       cache.Title,
				Description: cache.Description,
				ImgURL:      cache.ImgURL,
				BgColor:     cache.BgColor,
				MimeType:    cache.MimeType,
			}, nil
		}
	}

	return nil, fmt.Errorf("缓存未找到")
}

// fetchFaviconFromAPI 从外部API获取图标
// 使用 MonkNow API
func (s *FaviconService) fetchFaviconFromAPI(host string) (*FaviconInfo, error) {
	// 限流检查
	if !s.limiter.Allow() {
		log.Warn().Str("host", host).Msg("API 请求被限流")
		return nil, fmt.Errorf("请求过于频繁，请稍后重试")
	}

	// 移除端口号，只保留域名
	hostWithoutPort := s.removePort(host)

	// 构建 API URL
	apiURL := fmt.Sprintf("%s?url=%s", s.apiURL, url.QueryEscape(hostWithoutPort))

	// 创建请求
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置 User-Agent
	req.Header.Set("User-Agent", "DualTab/1.0 (https://github.com/yourusername/dualtab)")

	// 发起请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Warn().Err(err).Str("host", hostWithoutPort).Msg("请求API失败")
		return nil, fmt.Errorf("请求API失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		log.Warn().Int("status", resp.StatusCode).Str("host", hostWithoutPort).Msg("API返回错误状态")
		return nil, fmt.Errorf("API返回错误状态: %d", resp.StatusCode)
	}

	// 限制响应体大小（最大 1MB）
	limitedBody := io.LimitReader(resp.Body, 1<<20)

	// 读取响应体
	body, err := io.ReadAll(limitedBody)
	if err != nil {
		log.Warn().Err(err).Str("host", hostWithoutPort).Msg("读取响应失败")
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析 JSON 响应
	var apiResp MonkNowAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Warn().Err(err).Str("host", hostWithoutPort).Msg("解析响应失败")
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查响应状态
	if apiResp.Msg != "success" {
		log.Warn().Str("msg", apiResp.Msg).Str("host", hostWithoutPort).Msg("API返回错误")
		return nil, fmt.Errorf("API返回错误: %s", apiResp.Msg)
	}

	// 提取图标信息
	icon := apiResp.Data.Icon
	log.Info().Str("host", hostWithoutPort).Str("title", icon.Title).Msg("成功从API获取图标")

	return &FaviconInfo{
		Title:       icon.Title,
		Description: icon.Description,
		ImgURL:      icon.ImgURL,
		BgColor:     icon.BgColor,
		MimeType:    icon.MimeType,
	}, nil
}

// saveFaviconToCache 保存图标到缓存
func (s *FaviconService) saveFaviconToCache(host string, favicon *FaviconInfo) {
	// 移除端口号
	hostWithoutPort := s.removePort(host)

	cache := &model.FaviconCache{
		Host:        hostWithoutPort,
		Title:       favicon.Title,
		Description: favicon.Description,
		ImgURL:      favicon.ImgURL,
		BgColor:     favicon.BgColor,
		MimeType:    favicon.MimeType,
	}

	// 尝试保存，如果失败记录警告
	if err := s.cacheRepo.Create(cache); err != nil {
		log.Warn().Err(err).Str("host", hostWithoutPort).Msg("保存图标缓存失败")
	} else {
		log.Info().Str("host", hostWithoutPort).Msg("图标已保存到缓存")
	}
}

// removePort 移除域名中的端口号
func (s *FaviconService) removePort(host string) string {
	if idx := strings.Index(host, ":"); idx != -1 {
		return host[:idx]
	}
	return host
}
