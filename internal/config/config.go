package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config结构体用于存储整个项目的配置信息
type Config struct {
	DatabaseConfig DatabaseConfig
	ServerPort     int
}

// DatabaseConfig结构体用于存储数据库连接配置信息
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// LoadConfig函数用于从环境变量或.env文件中加载所有配置信息
func LoadConfig() (*Config, error) {
	// 尝试加载.env文件中的环境变量
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading.env file")
	}

	// 加载数据库配置
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	// 加载服务器端口配置
	serverPort, err := loadServerPort()
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseConfig: *dbConfig,
		ServerPort:     serverPort,
	}, nil
}

// loadDatabaseConfig函数用于从环境变量中加载数据库配置信息
func loadDatabaseConfig() (*DatabaseConfig, error) {
	return &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     parseInt(os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}, nil
}

// loadServerPort函数用于从环境变量中加载服务器端口配置信息
func loadServerPort() (int, error) {
	portStr := os.Getenv("SERVER_PORT")
	return parseInt(portStr), nil
}

// parseInt函数用于将字符串转换为整数
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("Invalid integer value: %s", s))
	}
	return i
}
