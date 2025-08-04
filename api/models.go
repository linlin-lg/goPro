package api

import (
	"time"
)

// User 用户模型
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Phone     string    `json:"phone"`
	Avatar    string    `json:"avatar"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DataItem 数据项模型
type DataItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Content     string    `json:"content"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SystemInfo 系统信息模型
type SystemInfo struct {
	Service   string `json:"service"`
	Version   string `json:"version"`
	Uptime    int64  `json:"uptime"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

// SystemStats 系统统计模型
type SystemStats struct {
	TotalRequests int     `json:"total_requests"`
	ActiveUsers   int     `json:"active_users"`
	DataCount     int     `json:"data_count"`
	SystemLoad    float64 `json:"system_load"`
	MemoryUsage   string  `json:"memory_usage"`
	CPUUsage      string  `json:"cpu_usage"`
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

// PaginationResponse 分页响应格式
type PaginationResponse struct {
	List  interface{} `json:"list"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int         `json:"total"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Avatar string `json:"avatar"`
}

// CreateDataRequest 创建数据请求
type CreateDataRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
}

// UpdateDataRequest 更新数据请求
type UpdateDataRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status"`
} 