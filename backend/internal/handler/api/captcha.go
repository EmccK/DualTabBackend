package api

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	store sync.Map // 简单的内存存储
}

// NewCaptchaHandler 创建验证码处理器
func NewCaptchaHandler() *CaptchaHandler {
	return &CaptchaHandler{}
}

// captchaData 验证码数据
type captchaData struct {
	code      string
	expiresAt time.Time
}

// generateCode 生成随机验证码
func generateCode(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果加密随机数生成失败，使用时间戳作为后备方案
			n = big.NewInt(int64(time.Now().UnixNano() % int64(len(charset))))
		}
		code[i] = charset[n.Int64()]
	}
	return string(code)
}

// generateCaptchaImage 生成简单的验证码图片
func generateCaptchaImage(code string, width, height int) ([]byte, error) {
	// 创建图片
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景色
	bgColor := color.RGBA{240, 240, 240, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// 简单绘制文字（使用像素点阵）
	textColor := color.RGBA{50, 50, 50, 255}
	charWidth := width / (len(code) + 1)
	startX := charWidth / 2

	for i, c := range code {
		drawChar(img, string(c), startX+i*charWidth, height/3, textColor)
	}

	// 添加干扰线
	lineColor := color.RGBA{180, 180, 180, 255}
	for i := 0; i < 3; i++ {
		y1, err := rand.Int(rand.Reader, big.NewInt(int64(height)))
		if err != nil {
			y1 = big.NewInt(int64(time.Now().UnixNano() % int64(height)))
		}
		y2, err := rand.Int(rand.Reader, big.NewInt(int64(height)))
		if err != nil {
			y2 = big.NewInt(int64(time.Now().UnixNano() % int64(height)))
		}
		drawLine(img, 0, int(y1.Int64()), width, int(y2.Int64()), lineColor)
	}

	// 编码为 PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// drawChar 简单绘制字符（5x7 点阵）
func drawChar(img *image.RGBA, char string, x, y int, c color.Color) {
	// 简化的字符点阵（仅支持数字和部分字母）
	patterns := map[string][]string{
		"0": {"01110", "10001", "10001", "10001", "10001", "10001", "01110"},
		"1": {"00100", "01100", "00100", "00100", "00100", "00100", "01110"},
		"2": {"01110", "10001", "00001", "00110", "01000", "10000", "11111"},
		"3": {"01110", "10001", "00001", "00110", "00001", "10001", "01110"},
		"4": {"00010", "00110", "01010", "10010", "11111", "00010", "00010"},
		"5": {"11111", "10000", "11110", "00001", "00001", "10001", "01110"},
		"6": {"01110", "10000", "11110", "10001", "10001", "10001", "01110"},
		"7": {"11111", "00001", "00010", "00100", "01000", "01000", "01000"},
		"8": {"01110", "10001", "10001", "01110", "10001", "10001", "01110"},
		"9": {"01110", "10001", "10001", "01111", "00001", "00001", "01110"},
		"A": {"01110", "10001", "10001", "11111", "10001", "10001", "10001"},
		"B": {"11110", "10001", "10001", "11110", "10001", "10001", "11110"},
		"C": {"01110", "10001", "10000", "10000", "10000", "10001", "01110"},
		"D": {"11110", "10001", "10001", "10001", "10001", "10001", "11110"},
		"E": {"11111", "10000", "10000", "11110", "10000", "10000", "11111"},
		"F": {"11111", "10000", "10000", "11110", "10000", "10000", "10000"},
	}

	pattern, ok := patterns[char]
	if !ok {
		return
	}

	scale := 3 // 放大倍数
	for row, line := range pattern {
		for col, bit := range line {
			if bit == '1' {
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						img.Set(x+col*scale+dx, y+row*scale+dy, c)
					}
				}
			}
		}
	}
}

// drawLine 绘制直线
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := 1, 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x1, y1, c)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// GetCaptcha 获取验证码
// GET /home/captcha?len=4&width=170&height=62&key=xxx
func (h *CaptchaHandler) GetCaptcha(c *gin.Context) {
	// 获取参数
	length, _ := strconv.Atoi(c.DefaultQuery("len", "4"))
	width, _ := strconv.Atoi(c.DefaultQuery("width", "170"))
	height, _ := strconv.Atoi(c.DefaultQuery("height", "62"))
	key := c.Query("key")

	if length < 4 {
		length = 4
	}
	if length > 6 {
		length = 6
	}

	// 生成验证码
	code := generateCode(length)

	// 存储验证码（5分钟有效）
	if key != "" {
		h.store.Store(key, captchaData{
			code:      code,
			expiresAt: time.Now().Add(5 * time.Minute),
		})
	}

	// 生成图片
	imgData, err := generateCaptchaImage(code, width, height)
	if err != nil {
		c.JSON(500, gin.H{"msg": "生成验证码失败"})
		return
	}

	// 返回 base64 图片
	c.Header("Content-Type", "image/png")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Data(200, "image/png", imgData)
}

// VerifyCaptcha 验证验证码
func (h *CaptchaHandler) VerifyCaptcha(key, code string) bool {
	if key == "" || code == "" {
		return false
	}

	value, ok := h.store.Load(key)
	if !ok {
		return false
	}

	data := value.(captchaData)

	// 检查是否过期
	if time.Now().After(data.expiresAt) {
		h.store.Delete(key)
		return false
	}

	// 验证码不区分大小写
	if data.code != code {
		return false
	}

	// 验证成功后删除
	h.store.Delete(key)
	return true
}

// GetCaptchaBase64 获取 Base64 格式的验证码（用于 API 返回）
func (h *CaptchaHandler) GetCaptchaBase64(length, width, height int, key string) (string, error) {
	code := generateCode(length)

	// 存储验证码
	if key != "" {
		h.store.Store(key, captchaData{
			code:      code,
			expiresAt: time.Now().Add(5 * time.Minute),
		})
	}

	imgData, err := generateCaptchaImage(code, width, height)
	if err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgData), nil
}
