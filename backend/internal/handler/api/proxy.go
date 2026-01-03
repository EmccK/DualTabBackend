package api

import (
	"dualtab-backend/internal/model"
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/response"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// ProxyHandler 代理服务处理器
type ProxyHandler struct {
	configRepo *repository.SystemConfigRepo
	client     *http.Client
}

// NewProxyHandler 创建代理服务处理器
func NewProxyHandler(configRepo *repository.SystemConfigRepo) *ProxyHandler {
	return &ProxyHandler{
		configRepo: configRepo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SearchSuggest 搜索建议代理
// GET /proxy/search-suggest?q=关键词
func (h *ProxyHandler) SearchSuggest(c *gin.Context) {
	// 检查是否启用搜索建议代理
	enabled := h.configRepo.GetValue(model.ConfigKeySearchSuggestOn)
	if enabled != "true" && enabled != "1" {
		response.BadRequest(c, "搜索建议代理未启用")
		return
	}

	query := c.Query("q")
	if query == "" {
		response.Success(c, []string{})
		return
	}

	// 代理 Google 搜索建议
	apiURL := fmt.Sprintf(
		"https://suggestqueries.google.com/complete/search?client=chrome&q=%s",
		url.QueryEscape(query),
	)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		response.InternalError(c, "创建请求失败")
		return
	}

	// 设置请求头，模拟浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := h.client.Do(req)
	if err != nil {
		response.InternalError(c, "请求搜索建议失败")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.InternalError(c, "读取响应失败")
		return
	}

	// Google 返回格式: ["query", ["suggestion1", "suggestion2", ...]]
	var result []interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		response.InternalError(c, "解析响应失败")
		return
	}

	// 提取建议列表
	if len(result) >= 2 {
		if suggestions, ok := result[1].([]interface{}); ok {
			strSuggestions := make([]string, len(suggestions))
			for i, s := range suggestions {
				if str, ok := s.(string); ok {
					strSuggestions[i] = str
				}
			}
			response.Success(c, strSuggestions)
			return
		}
	}

	response.Success(c, []string{})
}
