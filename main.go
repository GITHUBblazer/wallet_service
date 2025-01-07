package main

import (
	"fmt"
	"log"
	"net/http"

	"wallet-service/internal/api"
	"wallet-service/internal/config"
	"wallet-service/internal/database"
	"wallet-service/internal/logger"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	if cfg == nil {
		log.Fatal("配置结构体为nil，请检查配置加载逻辑")
	}
	log.Printf("Loaded config: %+v", cfg)

	db, err := database.ConnectDB(cfg.DatabaseConfig)
	if err != nil {
		logger.Log.Errorf("连接数据库失败: %v", err)
		if db == nil {
			logger.Log.Errorf("数据库连接对象为nil，具体错误: %v，请检查数据库连接逻辑", err)
		}
		return
	}
	defer db.Close()

	// 创建存储库和服务实例（可能会用到数据库连接等配置）
	repo := repository.NewRepository(db)
	if repo == nil {
		logger.Log.Errorf("存储库实例为nil，请检查存储库创建逻辑")
		return
	}
	walletService := service.NewWalletService(repo)
	if walletService == nil {
		logger.Log.Errorf("钱包服务实例为nil，请检查服务创建逻辑")
		return
	}

	// 创建API实例
	api := api.NewAPI(walletService)
	if api == nil {
		logger.Log.Errorf("API实例为nil，请检查API创建逻辑")
		return
	}

	// 定义HTTP路由并启动服务器（使用cfg.ServerPort中的端口号）
	router := api.Routes()
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	logger.Log.Infof("服务器启动，监听地址: %s", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		logger.Log.Errorf("启动服务器失败: %v", err)
	}
}
