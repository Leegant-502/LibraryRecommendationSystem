package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"library/api"
)

// SetupRoutes 设置API路由
func SetupRoutes(bookHandler *api.BookHandler, behaviorHandler *api.BehaviorTrackingHandler) *gin.Engine {
	router := gin.Default()

	// 添加中间件
	setupMiddleware(router)

	// 健康检查和系统信息
	router.GET("/health", healthCheck)
	router.GET("/version", versionInfo)

	// API版本分组
	v1 := router.Group("/api/v1")
	{
		// 用户行为追踪
		behavior := v1.Group("/behavior")
		{
			behavior.POST("/track", behaviorHandler.TrackUserBehavior)
		}

		// 推荐系统
		recommendations := v1.Group("/recommendations")
		{
			recommendations.GET("/personal", behaviorHandler.GetUserRecommendations)
			recommendations.GET("/popular", behaviorHandler.GetBehaviorBasedPopularBooks)
			recommendations.GET("/similar", behaviorHandler.GetBehaviorBasedSimilarBooks)
		}
	}

	// 兼容旧版本API（标记为废弃）
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

// healthCheck 健康检查端点
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "library-recommendation-system",
	})
}

// versionInfo 版本信息端点
func versionInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":     "1.0.0",
		"build_time":  "2024-01-01",
		"go_version":  "1.21",
		"description": "基于用户行为的图书推荐系统",
	})
}
