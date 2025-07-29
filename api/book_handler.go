package api

import (
	"github.com/gin-gonic/gin"
	"library/internal/service"
	"net/http"
	"strconv"
)

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

// RecordView 记录图书浏览
func (h *BookHandler) RecordView(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Title  string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if err := h.bookService.RecordBookView(req.UserID, req.Title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "浏览记录已保存"})
}

// RecordClick 记录图书点击
func (h *BookHandler) RecordClick(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Title  string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if err := h.bookService.RecordBookClick(req.UserID, req.Title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "点击记录已保存"})
}

// RecordRead 记录图书阅读
func (h *BookHandler) RecordRead(c *gin.Context) {
	var req struct {
		UserID          string `json:"user_id" binding:"required"`
		Title           string `json:"title" binding:"required"`
		ReadTimeMinutes int    `json:"read_time_minutes" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if err := h.bookService.RecordBookRead(req.UserID, req.Title, req.ReadTimeMinutes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "阅读记录已保存"})
}

// RecordStayTime 记录图书页面停留时间
func (h *BookHandler) RecordStayTime(c *gin.Context) {
	var req struct {
		UserID          string `json:"user_id" binding:"required"`
		Title           string `json:"title" binding:"required"`
		StayTimeSeconds int    `json:"stay_time_seconds" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if err := h.bookService.RecordBookStayTime(req.UserID, req.Title, req.StayTimeSeconds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "停留时间记录已保存"})
}

// GetRecommendations 获取图书推荐
func (h *BookHandler) GetRecommendations(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id 是必需的"})
		return
	}

	limit := 10 // 默认推荐10本
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	recommendations, err := h.bookService.GetRecommendations(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
		"count":           len(recommendations),
		"message":         "基于用户行为的个性化推荐",
	})
}

// GetPopularBooks 获取热门图书
func (h *BookHandler) GetPopularBooks(c *gin.Context) {
	limit := 10 // 默认返回10本
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	books, err := h.bookService.GetPopularBooks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"popular_books": books,
		"count":         len(books),
		"message":       "基于用户行为统计的热门图书",
	})
}

// GetSimilarBooks 获取相似图书
func (h *BookHandler) GetSimilarBooks(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title 是必需的"})
		return
	}

	limit := 10 // 默认返回10本
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	books, err := h.bookService.GetSimilarBooks(title, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"similar_books": books,
		"count":         len(books),
		"message":       "基于用户行为分析的相似图书推荐",
	})
}
