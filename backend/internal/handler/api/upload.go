package api

import (
	"dualtab-backend/pkg/response"
	"dualtab-backend/pkg/upload"

	"github.com/gin-gonic/gin"
)

// UploadHandler 上传处理器
type UploadHandler struct {
	uploader *upload.Uploader
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(uploader *upload.Uploader) *UploadHandler {
	return &UploadHandler{uploader: uploader}
}

// UploadImage 上传图片
// POST /upload/image
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取上传类型
	uploadType := c.DefaultPostForm("type", "icon")

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择文件")
		return
	}

	var url string
	switch uploadType {
	case "wallpaper":
		url, err = h.uploader.Upload(file, "wallpapers")
	case "avatar":
		url, err = h.uploader.Upload(file, "avatars")
	default:
		url, err = h.uploader.UploadIcon(file)
	}

	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.MonkNowSuccess(c, gin.H{"url": url})
}
