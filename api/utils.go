package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Response 统一响应格式
func Response(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	Response(c, http.StatusOK, "success", data)
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, code int, msg string) {
	Response(c, code, msg, nil)
}

// BadRequestResponse 请求参数错误响应
func BadRequestResponse(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusBadRequest, msg)
}

// NotFoundResponse 资源不存在响应
func NotFoundResponse(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusNotFound, msg)
}

// InternalServerErrorResponse 服务器内部错误响应
func InternalServerErrorResponse(c *gin.Context, msg string) {
	ErrorResponse(c, http.StatusInternalServerError, msg)
}

// GenerateID 生成唯一ID
func GenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GenerateMD5 生成MD5哈希
func GenerateMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// ParseInt 安全解析整数
func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

// ParseFloat 安全解析浮点数
func ParseFloat(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}
	
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return val
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	
	// 简单的邮箱格式验证
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	
	return true
}

// ValidatePhone 验证手机号格式
func ValidatePhone(phone string) bool {
	if phone == "" {
		return false
	}
	
	// 简单的手机号格式验证（中国大陆）
	if len(phone) != 11 {
		return false
	}
	
	// 检查是否都是数字
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}
	
	return true
}

// GetClientIP 获取客户端IP地址
func GetClientIP(c *gin.Context) string {
	// 尝试从各种头部获取真实IP
	ip := c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}
	
	ip = c.GetHeader("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For可能包含多个IP，取第一个
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	return c.ClientIP()
}

// GetUserAgent 获取用户代理
func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

// GetRequestID 获取请求ID（从头部或生成新的）
func GetRequestID(c *gin.Context) string {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = GenerateID()
	}
	return requestID
}

// LogRequest 记录请求日志
func LogRequest(c *gin.Context, duration time.Duration) {
	requestID := GetRequestID(c)
	clientIP := GetClientIP(c)
	userAgent := GetUserAgent(c)
	
	fmt.Printf("[%s] %s %s %s %s %v\n",
		requestID,
		c.Request.Method,
		c.Request.URL.Path,
		clientIP,
		userAgent,
		duration,
	)
}

// ParseJSONBody 解析JSON请求体
func ParseJSONBody(c *gin.Context, v interface{}) error {
	if err := c.ShouldBindJSON(v); err != nil {
		return err
	}
	return nil
}

// ParseQueryParams 解析查询参数
func ParseQueryParams(c *gin.Context, v interface{}) error {
	if err := c.ShouldBindQuery(v); err != nil {
		return err
	}
	return nil
}

// ParseFormData 解析表单数据
func ParseFormData(c *gin.Context, v interface{}) error {
	if err := c.ShouldBind(v); err != nil {
		return err
	}
	return nil
}

// JSONMarshal 安全的JSON序列化
func JSONMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JSONUnmarshal 安全的JSON反序列化
func JSONUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseTime 解析时间字符串
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeStr)
}

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// TruncateString 截断字符串
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// ContainsString 检查字符串切片是否包含指定字符串
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveString 从字符串切片中移除指定字符串
func RemoveString(slice []string, item string) []string {
	result := make([]string, 0)
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
} 