package api

import (
	//"GolandPro/storage"
	//"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// APIServer API服务器结构
type APIServer struct {
	router *gin.Engine
	logger *zap.Logger
	port   string
}

// NewAPIServer 创建新的API服务器
func NewAPIServer(port string) *APIServer {
	server := &APIServer{
		router: gin.Default(),
		port:   port,
	}

	// 初始化日志
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		logger = zap.NewNop()
	}
	server.logger = logger

	// 设置中间件
	server.setupMiddleware()

	// 设置路由
	server.setupRoutes()

	return server
}

// setupMiddleware 设置中间件
func (s *APIServer) setupMiddleware() {
	// 日志中间件
	s.router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		s.logger.Info("API Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("ip", c.ClientIP()),
		)
	})

	// CORS中间件
	s.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// setupRoutes 设置路由
func (s *APIServer) setupRoutes() {
	// API版本组
	v1 := s.router.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", s.healthCheck)

		// 用户相关接口
		users := v1.Group("/users")
		{
			users.GET("/:id", s.getUser)
			users.POST("/", s.createUser)
			users.PUT("/:id", s.updateUser)
			users.DELETE("/:id", s.deleteUser)
		}

		// 数据相关接口
		data := v1.Group("/data")
		{
			data.GET("/list", s.getDataList)
			data.POST("/", s.createData)
			data.GET("/:id", s.getData)
			data.PUT("/:id", s.updateData)
			data.DELETE("/:id", s.deleteData)
		}

		// 系统相关接口
		system := v1.Group("/system")
		{
			system.GET("/info", s.getSystemInfo)
			system.GET("/stats", s.getSystemStats)
		}
	}
}

// Start 启动API服务器
func (s *APIServer) Start() error {
	s.logger.Info("Starting API server", zap.String("port", s.port))
	return s.router.Run(":" + s.port)
}

// healthCheck 健康检查接口
func (s *APIServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "GolandPro API",
		"version":   "1.0.0",
	})
}

// getUser 获取用户信息
func (s *APIServer) getUser(c *gin.Context) {
	userID := c.Param("id")

	// 这里可以连接数据库查询用户信息
	user := map[string]interface{}{
		"id":      userID,
		"name":    "用户" + userID,
		"email":   "user" + userID + "@example.com",
		"created": time.Now().Add(-24 * time.Hour).Unix(),
		"status":  "active",
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": user,
	})
}

// createUser 创建用户
func (s *APIServer) createUser(c *gin.Context) {
	var user struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid request data: " + err.Error(),
		})
		return
	}

	// 这里可以保存用户到数据库
	newUser := map[string]interface{}{
		"id":      fmt.Sprintf("%d", time.Now().Unix()),
		"name":    user.Name,
		"email":   user.Email,
		"created": time.Now().Unix(),
		"status":  "active",
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 201,
		"msg":  "User created successfully",
		"data": newUser,
	})
}

// updateUser 更新用户信息
func (s *APIServer) updateUser(c *gin.Context) {
	userID := c.Param("id")

	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid request data: " + err.Error(),
		})
		return
	}

	// 这里可以更新数据库中的用户信息
	updatedUser := map[string]interface{}{
		"id":      userID,
		"name":    user.Name,
		"email":   user.Email,
		"updated": time.Now().Unix(),
		"status":  "active",
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "User updated successfully",
		"data": updatedUser,
	})
}

// deleteUser 删除用户
func (s *APIServer) deleteUser(c *gin.Context) {
	userID := c.Param("id")

	// 这里可以从数据库中删除用户
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "User deleted successfully",
		"data": map[string]string{
			"id": userID,
		},
	})
}

// getDataList 获取数据列表
func (s *APIServer) getDataList(c *gin.Context) {
	// 获取查询参数
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// 这里可以从数据库查询数据列表
	dataList := []map[string]interface{}{
		{
			"id":    "1",
			"title": "数据项1",
			"desc":  "这是第一个数据项的描述",
			"type":  "text",
		},
		{
			"id":    "2",
			"title": "数据项2",
			"desc":  "这是第二个数据项的描述",
			"type":  "image",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list":  dataList,
			"page":  page,
			"limit": limit,
			"total": len(dataList),
		},
	})
}

// createData 创建数据
func (s *APIServer) createData(c *gin.Context) {
	var data struct {
		Title string `json:"title" binding:"required"`
		Desc  string `json:"desc"`
		Type  string `json:"type"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid request data: " + err.Error(),
		})
		return
	}

	newData := map[string]interface{}{
		"id":      fmt.Sprintf("%d", time.Now().Unix()),
		"title":   data.Title,
		"desc":    data.Desc,
		"type":    data.Type,
		"created": time.Now().Unix(),
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": 201,
		"msg":  "Data created successfully",
		"data": newData,
	})
}

// getData 获取单个数据
func (s *APIServer) getData(c *gin.Context) {
	dataID := c.Param("id")

	data := map[string]interface{}{
		"id":      dataID,
		"title":   "数据项" + dataID,
		"desc":    "这是数据项" + dataID + "的详细描述",
		"type":    "text",
		"created": time.Now().Add(-24 * time.Hour).Unix(),
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

// updateData 更新数据
func (s *APIServer) updateData(c *gin.Context) {
	dataID := c.Param("id")

	var data struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Type  string `json:"type"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid request data: " + err.Error(),
		})
		return
	}

	updatedData := map[string]interface{}{
		"id":      dataID,
		"title":   data.Title,
		"desc":    data.Desc,
		"type":    data.Type,
		"updated": time.Now().Unix(),
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Data updated successfully",
		"data": updatedData,
	})
}

// deleteData 删除数据
func (s *APIServer) deleteData(c *gin.Context) {
	dataID := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Data deleted successfully",
		"data": map[string]string{
			"id": dataID,
		},
	})
}

// getSystemInfo 获取系统信息
func (s *APIServer) getSystemInfo(c *gin.Context) {
	info := map[string]interface{}{
		"service":   "GolandPro API",
		"version":   "1.0.0",
		"uptime":    time.Now().Unix(),
		"status":    "running",
		"timestamp": time.Now().Unix(),
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": info,
	})
}

// getSystemStats 获取系统统计信息
func (s *APIServer) getSystemStats(c *gin.Context) {
	stats := map[string]interface{}{
		"total_requests": 1000,
		"active_users":   50,
		"data_count":     200,
		"system_load":    0.75,
		"memory_usage":   "512MB",
		"cpu_usage":      "25%",
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": stats,
	})
}
