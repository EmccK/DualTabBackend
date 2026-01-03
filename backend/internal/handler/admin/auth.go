package admin

import (
	"dualtab-backend/internal/model"
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/jwt"
	"dualtab-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userRepo  *repository.AdminUserRepo
	jwtSecret string
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(userRepo *repository.AdminUserRepo, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入用户名和密码")
		return
	}

	// 查找用户
	user, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.ID, user.Username, h.jwtSecret)
	if err != nil {
		response.InternalError(c, "生成 Token 失败")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	response.Success(c, gin.H{
		"id":       userID,
		"username": username,
	})
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入旧密码和新密码")
		return
	}

	userID, _ := c.Get("user_id")
	user, err := h.userRepo.FindByID(userID.(uint))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		response.BadRequest(c, "旧密码错误")
		return
	}

	// 生成新密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}

	user.PasswordHash = string(hash)
	if err := h.userRepo.Update(user); err != nil {
		response.InternalError(c, "更新密码失败")
		return
	}

	response.Success(c, nil)
}

// CreateAdminRequest 创建管理员请求
type CreateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateAdmin 创建管理员（仅用于初始化）
func (h *AuthHandler) CreateAdmin(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入用户名和密码")
		return
	}

	// 检查用户名是否已存在
	if _, err := h.userRepo.FindByUsername(req.Username); err == nil {
		response.BadRequest(c, "用户名已存在")
		return
	}

	// 生成密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}

	user := &model.AdminUser{
		Username:     req.Username,
		PasswordHash: string(hash),
	}

	if err := h.userRepo.Create(user); err != nil {
		response.InternalError(c, "创建用户失败")
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}
