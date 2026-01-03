package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Msg:  "success",
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, Response{
		Msg: msg,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, msg)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, msg)
}

// NotFound 404 错误
func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, msg)
}

// InternalError 500 错误
func InternalError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, msg)
}

// PageData 分页数据结构
type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// SuccessWithPage 带分页的成功响应
func SuccessWithPage(c *gin.Context, list interface{}, total int64, page, size int) {
	Success(c, PageData{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	})
}
