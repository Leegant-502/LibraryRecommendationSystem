package routes

import (
	"github.com/gin-gonic/gin"
	"library/api"
)

// SetupRoutes 设置API路由
func SetupRoutes(bookHandler *api.BookHandler, behaviorHandler *api.BehaviorTrackingHandler) *gin.Engine {
	router := gin.Default()

	// 用户行为记录路由（保留原有接口以兼容现有前端）
	router.POST("/books/view", bookHandler.RecordView)          // 记录图书浏览
	router.POST("/books/click", bookHandler.RecordClick)        // 记录图书点击
	router.POST("/books/read", bookHandler.RecordRead)          // 记录图书阅读
	router.POST("/books/stay-time", bookHandler.RecordStayTime) // 记录图书页面停留时间

	// 新的统一用户行为追踪API
	router.POST("/behavior/track", behaviorHandler.TrackUserBehavior) // 统一的用户行为追踪接口

	// 推荐相关路由（保留原有接口以兼容现有前端）
	router.GET("/recommendations", bookHandler.GetRecommendations) // 获取个性化推荐
	router.GET("/books/popular", bookHandler.GetPopularBooks)      // 获取热门图书
	router.GET("/books/similar", bookHandler.GetSimilarBooks)      // 获取相似图书

	// 新的基于用户行为的推荐API
	router.GET("/behavior/recommendations", behaviorHandler.GetUserRecommendations) // 基于用户行为的个性化推荐
	router.GET("/behavior/popular", behaviorHandler.GetBehaviorBasedPopularBooks)   // 基于用户行为的热门图书
	router.GET("/behavior/similar", behaviorHandler.GetBehaviorBasedSimilarBooks)   // 基于用户行为的相似图书

	return router
}
