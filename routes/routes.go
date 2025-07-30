package routes

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"library/api"
)

// SetupRoutes 设置API路由
func SetupRoutes(unifiedHandler *api.UnifiedHandler, bookHandler *api.BookHandler) *gin.Engine {
	router := gin.Default()

	// 添加中间件
	setupMiddleware(router)

	// 健康检查和系统信息
	router.GET("/health", unifiedHandler.HealthCheck)
	router.GET("/version", unifiedHandler.GetVersion)

	// API版本分组
	v1 := router.Group("/api/v1")
	{
		// 用户行为追踪
		v1.POST("/behavior/track", unifiedHandler.TrackUserBehavior)

		// 推荐系统
		recommendations := v1.Group("/recommendations")
		{
			recommendations.GET("/personal", unifiedHandler.GetPersonalizedRecommendations)
			recommendations.GET("/popular", unifiedHandler.GetPopularBooks)
			recommendations.GET("/similar", unifiedHandler.GetSimilarBooks)
		}
	}

	// 兼容旧版本API（标记为废弃，逐步迁移）
	deprecated := router.Group("/deprecated")
	{
		deprecated.POST("/books/view", bookHandler.RecordView)
		deprecated.POST("/books/click", bookHandler.RecordClick)
		deprecated.POST("/books/read", bookHandler.RecordRead)
		deprecated.POST("/books/stay-time", bookHandler.RecordStayTime)
		deprecated.GET("/recommendations", bookHandler.GetRecommendations)
		deprecated.GET("/books/popular", bookHandler.GetPopularBooks)
		deprecated.GET("/books/similar", bookHandler.GetSimilarBooks)
	}

	return router
}

// setupMiddleware 设置中间件
func setupMiddleware(router *gin.Engine) {
	// CORS中间件
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 请求日志中间件
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// 恢复中间件
	router.Use(gin.Recovery())
}
