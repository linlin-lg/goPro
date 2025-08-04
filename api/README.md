# GolandPro API 服务

这是一个为app提供简单接口能力的RESTful API服务。

## 功能特性

- 🚀 基于Gin框架的高性能HTTP服务
- 📝 完整的RESTful API接口
- 🔒 CORS跨域支持
- 📊 请求日志记录
- ⚙️ 灵活的配置管理
- 🛡️ 统一的错误处理
- 📋 数据验证和格式化

## 快速开始

### 1. 启动服务

```bash
# 使用默认配置启动
go run api/main.go

# 指定端口启动
go run api/main.go -port 9090

# 指定配置文件启动
go run api/main.go -config config.json
```

### 2. 环境变量配置

```bash
export API_PORT=8080
export API_HOST=0.0.0.0
export DB_HOST=localhost
export DB_PORT=3306
export DB_USERNAME=root
export DB_PASSWORD=password
export DB_DATABASE=golandpro
export LOG_LEVEL=info
```

## API 接口文档

### 基础信息

- **基础URL**: `http://localhost:8080/api/v1`
- **内容类型**: `application/json`
- **字符编码**: `UTF-8`

### 统一响应格式

```json
{
  "code": 200,
  "msg": "success",
  "data": {}
}
```

### 接口列表

#### 1. 健康检查

- **GET** `/api/v1/health`
- **描述**: 检查服务健康状态
- **响应示例**:
```json
{
  "status": "ok",
  "timestamp": 1640995200,
  "service": "GolandPro API",
  "version": "1.0.0"
}
```

#### 2. 用户管理

##### 获取用户信息
- **GET** `/api/v1/users/{id}`
- **参数**: `id` - 用户ID
- **响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": "123",
    "name": "张三",
    "email": "zhangsan@example.com",
    "created": 1640995200,
    "status": "active"
  }
}
```

##### 创建用户
- **POST** `/api/v1/users/`
- **请求体**:
```json
{
  "name": "张三",
  "email": "zhangsan@example.com"
}
```
- **响应示例**:
```json
{
  "code": 201,
  "msg": "User created successfully",
  "data": {
    "id": "1640995200",
    "name": "张三",
    "email": "zhangsan@example.com",
    "created": 1640995200,
    "status": "active"
  }
}
```

##### 更新用户
- **PUT** `/api/v1/users/{id}`
- **请求体**:
```json
{
  "name": "李四",
  "email": "lisi@example.com"
}
```

##### 删除用户
- **DELETE** `/api/v1/users/{id}`

#### 3. 数据管理

##### 获取数据列表
- **GET** `/api/v1/data/list`
- **查询参数**:
  - `page` - 页码 (默认: 1)
  - `limit` - 每页数量 (默认: 10)
- **响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": "1",
        "title": "数据项1",
        "desc": "这是第一个数据项的描述",
        "type": "text"
      }
    ],
    "page": 1,
    "limit": 10,
    "total": 1
  }
}
```

##### 创建数据
- **POST** `/api/v1/data/`
- **请求体**:
```json
{
  "title": "新数据项",
  "desc": "数据项描述",
  "type": "text"
}
```

##### 获取单个数据
- **GET** `/api/v1/data/{id}`

##### 更新数据
- **PUT** `/api/v1/data/{id}`

##### 删除数据
- **DELETE** `/api/v1/data/{id}`

#### 4. 系统信息

##### 获取系统信息
- **GET** `/api/v1/system/info`

##### 获取系统统计
- **GET** `/api/v1/system/stats`

## 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 配置说明

### 配置文件格式 (config.json)

```json
{
  "server": {
    "port": "8080",
    "host": "0.0.0.0",
    "read_timeout": 30,
    "write_timeout": 30,
    "idle_timeout": 60
  },
  "database": {
    "host": "localhost",
    "port": 3306,
    "username": "root",
    "password": "",
    "database": "golandpro",
    "charset": "utf8mb4"
  },
  "log": {
    "level": "info",
    "output": "stdout",
    "max_size": 100,
    "max_backups": 3,
    "max_age": 28,
    "compress": true
  },
  "cors": {
    "allow_origins": ["*"],
    "allow_methods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    "allow_headers": ["Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"],
    "expose_headers": ["Content-Length"],
    "allow_credentials": true,
    "max_age": 86400
  }
}
```

## 开发指南

### 添加新的API接口

1. 在 `server.go` 的 `setupRoutes()` 方法中添加路由
2. 实现对应的处理函数
3. 在 `models.go` 中定义相关的数据结构
4. 使用 `utils.go` 中的工具函数处理请求和响应

### 示例：添加新的接口

```go
// 在 setupRoutes() 中添加
v1.GET("/custom", s.customHandler)

// 实现处理函数
func (s *APIServer) customHandler(c *gin.Context) {
    // 处理逻辑
    SuccessResponse(c, gin.H{"message": "Custom API"})
}
```

## 部署说明

### Docker 部署

```dockerfile
FROM golang:1.17-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o api-server api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api-server .
EXPOSE 8080
CMD ["./api-server"]
```

### 系统服务部署

创建 systemd 服务文件 `/etc/systemd/system/golandpro-api.service`:

```ini
[Unit]
Description=GolandPro API Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/golandpro
ExecStart=/path/to/golandpro/api-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl enable golandpro-api
sudo systemctl start golandpro-api
```

## 许可证

MIT License 