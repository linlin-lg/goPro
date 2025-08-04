package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Config API服务配置
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Log      LogConfig      `json:"log"`
	Cors     CorsConfig     `json:"cors"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string `json:"port"`
	Host         string `json:"host"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	IdleTimeout  int    `json:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Charset  string `json:"charset"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level"`
	Output     string `json:"output"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// CorsConfig CORS配置
type CorsConfig struct {
	AllowOrigins     []string `json:"allow_origins"`
	AllowMethods     []string `json:"allow_methods"`
	AllowHeaders     []string `json:"allow_headers"`
	ExposeHeaders    []string `json:"expose_headers"`
	AllowCredentials bool     `json:"allow_credentials"`
	MaxAge           int      `json:"max_age"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         "8080",
			Host:         "0.0.0.0",
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  60,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Username: "root",
			Password: "",
			Database: "golandpro",
			Charset:  "utf8mb4",
		},
		Log: LogConfig{
			Level:      "info",
			Output:     "stdout",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		},
		Cors: CorsConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
			ExposeHeaders: []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge: 86400,
		},
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()
	
	if configPath == "" {
		configPath = "config.json"
	}
	
	file, err := os.Open(configPath)
	if err != nil {
		// 如果文件不存在，返回默认配置
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %v", err)
	}
	
	return config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, configPath string) error {
	if configPath == "" {
		configPath = "config.json"
	}
	
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode config: %v", err)
	}
	
	return nil
}

// GetEnvConfig 从环境变量获取配置
func GetEnvConfig() *Config {
	config := DefaultConfig()
	
	// 服务器配置
	if port := os.Getenv("API_PORT"); port != "" {
		config.Server.Port = port
	}
	if host := os.Getenv("API_HOST"); host != "" {
		config.Server.Host = host
	}
	if readTimeout := os.Getenv("API_READ_TIMEOUT"); readTimeout != "" {
		if val, err := strconv.Atoi(readTimeout); err == nil {
			config.Server.ReadTimeout = val
		}
	}
	if writeTimeout := os.Getenv("API_WRITE_TIMEOUT"); writeTimeout != "" {
		if val, err := strconv.Atoi(writeTimeout); err == nil {
			config.Server.WriteTimeout = val
		}
	}
	
	// 数据库配置
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if val, err := strconv.Atoi(dbPort); err == nil {
			config.Database.Port = val
		}
	}
	if dbUser := os.Getenv("DB_USERNAME"); dbUser != "" {
		config.Database.Username = dbUser
	}
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		config.Database.Password = dbPass
	}
	if dbName := os.Getenv("DB_DATABASE"); dbName != "" {
		config.Database.Database = dbName
	}
	
	// 日志配置
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.Log.Level = logLevel
	}
	if logOutput := os.Getenv("LOG_OUTPUT"); logOutput != "" {
		config.Log.Output = logOutput
	}
	
	return config
}

// GetDatabaseDSN 获取数据库连接字符串
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.Charset,
	)
}

// GetServerAddr 获取服务器地址
func (c *Config) GetServerAddr() string {
	return c.Server.Host + ":" + c.Server.Port
} 