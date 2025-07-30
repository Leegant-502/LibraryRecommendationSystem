package main

import (
	"context"
	"fmt"
	"library/internal/bookFetch"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// 配置数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库表
	err = db.AutoMigrate(&model.BookInfo{}, &model.UserBehavior{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化依赖
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo, cfg.Gorse.Endpoint, cfg.Gorse.APIKey)

	// 创建处理器
	unifiedHandler := api.NewUnifiedHandler(bookService)
	bookHandler := api.NewBookHandler(bookService) // 保留用于兼容性

	// 设置路由
	mux := routes.SetupRoutes(unifiedHandler, bookHandler)

	// 创建服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("启动图书数据定时同步任务...")
		bookFetch.TimelyFetchAndSaveBooks(db)
	}()

	// 启动服务器（非阻塞）
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 优雅关闭服务器，等待现有连接完成
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
