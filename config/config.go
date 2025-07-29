package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Gorse    GorseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type GorseConfig struct {
	Endpoint string
	APIKey   string
}

// getEnv 从环境变量获取值，如果环境变量不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadEnv 加载环境变量文件
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("警告: .env 文件未找到，将使用默认配置: %v", err)
	}
}

func NewConfig() *Config {
	// 加载环境变量
	LoadEnv()

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "library"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Gorse: GorseConfig{
			Endpoint: getEnv("GORSE_ENDPOINT", "http://localhost:8088"),
			APIKey:   getEnv("GORSE_API_KEY", ""),
		},
	}
}
