package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"library/internal/service"
)

// UnifiedHandler 统一的API处理器
// 专注于HTTP请求处理，业务逻辑委托给Service层
type UnifiedHandler struct {
	bookService service.BookServiceInterface
}

// NewUnifiedHandler 创建新的统一处理器
func NewUnifiedHandler(bookService service.BookServiceInterface) *UnifiedHandler {
	return &UnifiedHandler{bookService: bookService}
}

// TrackUserBehavior 用户行为追踪
func (h *UnifiedHandler) TrackUserBehavior(c *gin.Context) {
	var req service.UserBehaviorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "参数错误",
			"details": err.Error(),
		})
		return
	}

	if err := h.bookService.RecordUserBehavior(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "记录用户行为失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "用户行为记录成功",
		"behavior_type": req.BehaviorType,
		"user_id":       req.UserID,
		"book_title":    req.BookTitle,
	})
}

// GetPersonalizedRecommendations 获取个性化推荐
func (h *UnifiedHandler) GetPersonalizedRecommendations(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id 参数是必需的",
		})
		return
	}

	limit := 10 // 默认推荐10本
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	recommendations, err := h.bookService.GetRecommendations(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取推荐失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"recommendations": recommendations,
		"count":           len(recommendations),
		"user_id":         userID,
		"algorithm":       "基于用户行为的协同过滤推荐",
	})
}

// GetPopularBooks 获取热门图书
func (h *UnifiedHandler) GetPopularBooks(c *gin.Context) {
	limit := 10 // 默认返回10本
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	books, err := h.bookService.GetPopularBooks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取热门图书失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"popular_books": books,
		"count":         len(books),
		"algorithm":     "基于用户行为统计的热门度排序",
	})
}

// GetSimilarBooks 获取相似图书
func (h *UnifiedHandler) GetSimilarBooks(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "title 参数是必需的",
		})
		return
	}

	limit := 10 // 默认返回10本
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	books, err := h.bookService.GetSimilarBooks(title, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取相似图书失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"similar_books": books,
		"count":         len(books),
		"base_title":    title,
		"algorithm":     "基于用户行为的物品协同过滤",
	})
}

// HealthCheck 健康检查
func (h *UnifiedHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "library-recommendation-system",
		"version": "1.0.0",
	})
}

// GetVersion 获取版本信息
func (h *UnifiedHandler) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":     "1.0.0",
		"build_time":  "2024-01-01",
		"go_version":  "1.21",
		"description": "基于用户行为的图书推荐系统",
		"features": []string{
			"用户行为追踪",
			"个性化推荐",
			"热门图书推荐",
			"相似图书推荐",
		},
	})
}
