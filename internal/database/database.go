package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"wallet-service/internal/config"
)

// ConnectDB ConnectDB函数接受数据库配置结构体并返回一个数据库连接对象和可能的错误
func ConnectDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
