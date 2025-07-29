package api

import (
	"github.com/gin-gonic/gin"
	"library/internal/service"
	"net/http"
	"strconv"
)

// BehaviorTrackingHandler 用户行为追踪处理器
type BehaviorTrackingHandler struct {
	bookService *service.BookService
}

// NewBehaviorTrackingHandler 创建新的行为追踪处理器
func NewBehaviorTrackingHandler(bookService *service.BookService) *BehaviorTrackingHandler {
	return &BehaviorTrackingHandler{bookService: bookService}
}

// TrackUserBehavior 统一的用户行为追踪接口
func (h *BehaviorTrackingHandler) TrackUserBehavior(c *gin.Context) {
	var req struct {
		UserID          string                 `json:"user_id" binding:"required"`
		BookTitle       string                 `json:"book_title" binding:"required"`
		BehaviorType    string                 `json:"behavior_type" binding:"required"` // "view", "click", "read", "stay_time"
		StayTimeSeconds *int                   `json:"stay_time_seconds,omitempty"`      // 停留时间（秒）
		ReadTimeMinutes *int                   `json:"read_time_minutes,omitempty"`      // 阅读时间（分钟）
		Extra           map[string]interface{} `json:"extra,omitempty"`                  // 额外信息
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	var err error
	switch req.BehaviorType {
	case "view":
		err = h.bookService.RecordBookView(req.UserID, req.BookTitle)
	case "click":
		err = h.bookService.RecordBookClick(req.UserID, req.BookTitle)
	case "read":
		if req.ReadTimeMinutes == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "阅读行为需要提供 read_time_minutes 参数"})
			return
		}
		err = h.bookService.RecordBookRead(req.UserID, req.BookTitle, *req.ReadTimeMinutes)
	case "stay_time":
		if req.StayTimeSeconds == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "停留时间行为需要提供 stay_time_seconds 参数"})
			return
		}
		err = h.bookService.RecordBookStayTime(req.UserID, req.BookTitle, *req.StayTimeSeconds)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的行为类型: " + req.BehaviorType})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "用户行为记录成功",
		"behavior_type": req.BehaviorType,
		"user_id":       req.UserID,
		"book_title":    req.BookTitle,
	})
}

// GetUserRecommendations 获取基于用户行为的推荐
func (h *BehaviorTrackingHandler) GetUserRecommendations(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id 是必需的"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"recommendations": recommendations,
		"count":           len(recommendations),
		"user_id":         userID,
		"algorithm":       "基于用户行为的协同过滤推荐",
		"message":         "推荐结果基于您的浏览、点击和停留时间等行为数据生成",
	})
}

// GetBehaviorBasedPopularBooks 获取基于用户行为的热门图书
func (h *BehaviorTrackingHandler) GetBehaviorBasedPopularBooks(c *gin.Context) {
	limit := 10 // 默认返回10本
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	books, err := h.bookService.GetPopularBooks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"popular_books": books,
		"count":         len(books),
		"algorithm":     "基于用户行为统计的热门度排序",
		"message":       "热门图书基于所有用户的点击、浏览和停留时间等行为数据统计生成",
	})
}

// GetBehaviorBasedSimilarBooks 获取基于用户行为的相似图书
func (h *BehaviorTrackingHandler) GetBehaviorBasedSimilarBooks(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title 是必需的"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"similar_books": books,
		"count":         len(books),
		"base_title":    title,
		"algorithm":     "基于用户行为的物品协同过滤",
		"message":       "相似图书基于用户对图书的行为模式相似性推荐",
	})
}
