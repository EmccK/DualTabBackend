package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// Uploader 文件上传器
type Uploader struct {
	BasePath string // 存储根目录
	BaseURL  string // 访问 URL 前缀
}

// NewUploader 创建上传器
func NewUploader(basePath, baseURL string) *Uploader {
	return &Uploader{
		BasePath: basePath,
		BaseURL:  baseURL,
	}
}

// AllowedImageTypes 允许的图片类型
var AllowedImageTypes = map[string]bool{
	"image/png":     true,
	"image/jpeg":    true,
	"image/jpg":     true,
	"image/svg+xml": true,
	"image/webp":    true,
	"image/gif":     true,
}

// MaxFileSize 最大文件大小 (2MB)
const MaxFileSize = 2 * 1024 * 1024

// UploadIcon 上传图标
func (u *Uploader) UploadIcon(file *multipart.FileHeader) (string, error) {
	// 验证文件类型
	contentType := file.Header.Get("Content-Type")
	if !AllowedImageTypes[contentType] {
		return "", fmt.Errorf("不支持的文件类型: %s", contentType)
	}

	// 验证文件大小
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("文件大小不能超过 2MB")
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// 根据 MIME 类型推断扩展名
		switch contentType {
		case "image/png":
			ext = ".png"
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/svg+xml":
			ext = ".svg"
		case "image/webp":
			ext = ".webp"
		case "image/gif":
			ext = ".gif"
		default:
			ext = ".png"
		}
	}
	filename := uuid.New().String() + strings.ToLower(ext)

	// 确保目录存在
	iconDir := filepath.Join(u.BasePath, "icons")
	if err := os.MkdirAll(iconDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 保存文件
	dst := filepath.Join(iconDir, filename)
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	// 返回访问 URL
	return fmt.Sprintf("%s/icons/%s", u.BaseURL, filename), nil
}
