package service

import (
	"dualtab-backend/internal/repository"
	"fmt"
	"net/url"
	"strings"
)

// FaviconService 网站图标服务
type FaviconService struct {
	cacheRepo *repository.FaviconCacheRepo
}

// NewFaviconService 创建网站图标服务
func NewFaviconService(cacheRepo *repository.FaviconCacheRepo) *FaviconService {
	return &FaviconService{
		cacheRepo: cacheRepo,
	}
}

// FaviconInfo 图标信息
type FaviconInfo struct {
	Title    string
	ImgURL   string
	BgColor  string
	MimeType string
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
		return favicon, nil
	}

	// 缓存未命中，返回错误
	return nil, fmt.Errorf("未找到图标缓存")
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
			Title:    cache.Title,
			ImgURL:   cache.ImgURL,
			BgColor:  cache.BgColor,
			MimeType: cache.MimeType,
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
				Title:    cache.Title,
				ImgURL:   cache.ImgURL,
				BgColor:  cache.BgColor,
				MimeType: cache.MimeType,
			}, nil
		}
	}

	return nil, fmt.Errorf("缓存未找到")
}
