package main

import (
	"fmt"
	"log"
	"net/http"

	"library/api"
	"library/config"
	"library/internal/model"
	"library/internal/repository"
	"library/internal/service"
	"library/routes"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg := config.NewConfig()

	// 构建 PostgreSQL 连接字符串
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// 初始化数据库连接
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(&model.BookInfo{}, &model.UserBehavior{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化依赖
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo, cfg.Gorse.Endpoint, cfg.Gorse.APIKey)
	bookHandler := api.NewBookHandler(bookService)
	behaviorHandler := api.NewBehaviorTrackingHandler(bookService)

	// 设置路由
	mux := routes.SetupRoutes(bookHandler, behaviorHandler)

	// 创建服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: mux,
	}

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
