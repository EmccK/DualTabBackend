package api

import (
	"crypto/rand"
	"dualtab-backend/internal/model"
	"dualtab-backend/internal/repository"
	"dualtab-backend/pkg/response"
	"encoding/hex"
	"fmt"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler 用户 API 处理器
type UserHandler struct {
	userRepo     *repository.UserRepo
	userDataRepo *repository.UserDataRepo
}

// NewUserHandler 创建用户 API 处理器
func NewUserHandler(userRepo *repository.UserRepo, userDataRepo *repository.UserDataRepo) *UserHandler {
	return &UserHandler{
		userRepo:     userRepo,
		userDataRepo: userDataRepo,
	}
}

// generateSecret 生成随机 secret
func generateSecret() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
// POST /user/login
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请输入有效的邮箱和密码")
		return
	}

	// 查找用户
	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		response.Error(c, 401, "邮箱或密码错误")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, 401, "邮箱或密码错误")
		return
	}

	// 更新最后访问时间
	now := time.Now().Unix()
	h.userRepo.UpdateLastVisit(user.ID, now)
	user.LastVisitAt = now

	response.MonkNowSuccess(c, user.ToResponse())
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Captcha  string `json:"captcha"` // 验证码（可选）
	Key      string `json:"key"`     // 验证码 key（可选）
}

// validatePassword 验证密码强度
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("密码长度至少8位")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return fmt.Errorf("密码必须包含大写字母、小写字母和数字")
	}

	return nil
}

// Register 用户注册
// POST /user/register
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请输入有效的邮箱和密码（密码至少8位）")
		return
	}

	// 验证密码强度
	if err := validatePassword(req.Password); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	// 检查邮箱是否已存在
	if _, err := h.userRepo.FindByEmail(req.Email); err == nil {
		response.Error(c, 400, "该邮箱已被注册")
		return
	}

	// 生成密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, 500, "注册失败")
		return
	}

	// 创建用户
	user := &model.User{
		Email:        req.Email,
		Name:         req.Email, // 默认使用邮箱作为昵称
		PasswordHash: string(hash),
		Secret:       generateSecret(),
		IsActivate:   1,
		LastVisitAt:  time.Now().Unix(),
	}

	if err := h.userRepo.Create(user); err != nil {
		response.Error(c, 500, "注册失败")
		return
	}

	response.MonkNowSuccess(c, gin.H{})
}

// GetUserData 获取用户数据
// GET /user/data/info?type=icons
func (h *UserHandler) GetUserData(c *gin.Context) {
	// 从请求头获取 secret
	secret := c.GetHeader("secret")
	if secret == "" {
		response.Error(c, 401, "未授权")
		return
	}

	// 查找用户
	user, err := h.userRepo.FindBySecret(secret)
	if err != nil {
		response.Error(c, 401, "无效的授权")
		return
	}

	dataType := c.Query("type")
	if dataType == "" {
		response.Error(c, 400, "请指定数据类型")
		return
	}

	// 获取用户数据
	data, err := h.userDataRepo.GetByType(user.ID, dataType)
	if err != nil {
		// 没有数据时返回空
		response.MonkNowSuccess(c, gin.H{"data": nil})
		return
	}

	response.MonkNowSuccess(c, gin.H{"data": data.Data})
}

// UpdateUserDataRequest 更新用户数据请求
type UpdateUserDataRequest struct {
	Type string `json:"type" binding:"required"`
	Data string `json:"data" binding:"required"`
}

// UpdateUserData 更新用户数据
// PUT /user/data/update
func (h *UserHandler) UpdateUserData(c *gin.Context) {
	// 从请求头获取 secret
	secret := c.GetHeader("secret")
	if secret == "" {
		response.Error(c, 401, "未授权")
		return
	}

	// 查找用户
	user, err := h.userRepo.FindBySecret(secret)
	if err != nil {
		response.Error(c, 401, "无效的授权")
		return
	}

	var req UpdateUserDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	// 验证数据类型
	validTypes := map[string]bool{
		"icons":      true,
		"common":     true,
		"background": true,
		"searcher":   true,
		"sidebar":    true,
		"todos":      true,
		"standby":    true,
	}
	if !validTypes[req.Type] {
		response.Error(c, 400, "无效的数据类型")
		return
	}

	// 更新数据
	if err := h.userDataRepo.Upsert(user.ID, req.Type, req.Data); err != nil {
		response.Error(c, 500, "更新数据失败")
		return
	}

	// 标记用户已备份
	if user.IsBackup == 0 {
		user.IsBackup = 1
		h.userRepo.Update(user)
	}

	response.MonkNowSuccess(c, gin.H{})
}

// ChangeNameRequest 修改昵称请求
type ChangeNameRequest struct {
	Name string `json:"name" binding:"required"`
}

// ChangeName 修改昵称
// PUT /user/changename
func (h *UserHandler) ChangeName(c *gin.Context) {
	secret := c.GetHeader("secret")
	if secret == "" {
		response.Error(c, 401, "未授权")
		return
	}

	user, err := h.userRepo.FindBySecret(secret)
	if err != nil {
		response.Error(c, 401, "无效的授权")
		return
	}

	var req ChangeNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请输入昵称")
		return
	}

	user.Name = req.Name
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, 500, "更新失败")
		return
	}

	response.MonkNowSuccess(c, gin.H{})
}

// ChangeAvatarRequest 修改头像请求
type ChangeAvatarRequest struct {
	Avatar string `json:"avatar" binding:"required"`
}

// ChangeAvatar 修改头像
// PUT /user/changeavatar
func (h *UserHandler) ChangeAvatar(c *gin.Context) {
	secret := c.GetHeader("secret")
	if secret == "" {
		response.Error(c, 401, "未授权")
		return
	}

	user, err := h.userRepo.FindBySecret(secret)
	if err != nil {
		response.Error(c, 401, "无效的授权")
		return
	}

	var req ChangeAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请提供头像 URL")
		return
	}

	user.Avatar = req.Avatar
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, 500, "更新失败")
		return
	}

	response.MonkNowSuccess(c, gin.H{})
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	Old string `json:"old" binding:"required"`
	New string `json:"new" binding:"required,min=8"`
}

// ChangePassword 修改密码
// PUT /user/changepwd
func (h *UserHandler) ChangePassword(c *gin.Context) {
	secret := c.GetHeader("secret")
	if secret == "" {
		response.Error(c, 401, "未授权")
		return
	}

	user, err := h.userRepo.FindBySecret(secret)
	if err != nil {
		response.Error(c, 401, "无效的授权")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请输入旧密码和新密码（新密码至少8位）")
		return
	}

	// 验证新密码强度
	if err := validatePassword(req.New); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Old)); err != nil {
		response.Error(c, 400, "旧密码错误")
		return
	}

	// 生成新密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.New), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, 500, "修改密码失败")
		return
	}

	user.PasswordHash = string(hash)
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, 500, "修改密码失败")
		return
	}

	response.MonkNowSuccess(c, gin.H{})
}
