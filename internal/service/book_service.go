package service

import (
	"fmt"
	"library/internal/domain/model"
	"library/internal/gorse"
	"library/internal/repository"
	"time"
)

// BookService 处理图书相关的业务逻辑
type BookService struct {
	bookRepo    repository.BookRepository
	gorseClient *gorse.Client
}

// NewBookService 创建新的 BookService 实例
func NewBookService(bookRepo repository.BookRepository, gorseEndpoint, gorseAPIKey string) *BookService {
	return &BookService{
		bookRepo:    bookRepo,
		gorseClient: gorse.NewClient(gorseEndpoint, gorseAPIKey),
	}
}

// RecordBookView 记录图书浏览行为
func (s *BookService) RecordBookView(userID, title string) error {
	return s.gorseClient.InsertFeedback("view", userID, title, time.Now().Unix(), nil)
}

// RecordBookClick 记录图书点击行为
func (s *BookService) RecordBookClick(userID, title string) error {
	return s.gorseClient.InsertFeedback("click", userID, title, time.Now().Unix(), nil)
}

// RecordBookRead 记录图书阅读行为
func (s *BookService) RecordBookRead(userID, title string, readTimeMinutes int) error {
	extra := map[string]interface{}{
		"read_time": readTimeMinutes,
	}
	return s.gorseClient.InsertFeedback("read", userID, title, time.Now().Unix(), extra)
}

// RecordBookStayTime 记录图书页面停留时间
func (s *BookService) RecordBookStayTime(userID, title string, stayTimeSeconds int) error {
	extra := map[string]interface{}{
		"stay_time": stayTimeSeconds,
	}

	// 根据停留时间决定反馈类型
	feedbackType := "view"
	if stayTimeSeconds >= 30 { // 停留超过30秒视为深度阅读
		feedbackType = "read"
	}

	return s.gorseClient.InsertFeedback(feedbackType, userID, title, time.Now().Unix(), extra)
}

// GetRecommendations 获取图书推荐，包含对新用户的处理
func (s *BookService) GetRecommendations(userID string, limit int) ([]*model.BookInfo, error) {
	// 先尝试获取个性化推荐
	recommendationTitles, err := s.gorseClient.GetRecommend(userID, "", limit)
	if err != nil {
		return nil, fmt.Errorf("获取推荐失败: %v", err)
	}

	// 如果没有个性化推荐结果，使用默认推荐策略
	if len(recommendationTitles) == 0 {
		recommendationTitles, err = s.getDefaultRecommendations(limit)
		if err != nil {
			return nil, err
		}
	}

	// 根据标题获取完整的图书信息
	return s.getBooksByTitles(recommendationTitles)
}

// getDefaultRecommendations 获取默认推荐（针对新用户）
func (s *BookService) getDefaultRecommendations(limit int) ([]string, error) {
	// 策略1：获取热门图书（占比60%）
	popularLimit := int(float64(limit) * 0.6)
	popularBooks, err := s.gorseClient.GetPopular("", popularLimit)
	if err != nil {
		return nil, fmt.Errorf("获取热门图书失败: %v", err)
	}

	// 策略2：获取最新图书（占比40%）
	latestLimit := limit - len(popularBooks)
	latestBooks, err := s.gorseClient.GetLatest("", latestLimit)
	if err != nil {
		latestBooks = []string{} // 如果获取最新图书失败，使用空列表
	}

	// 合并推荐结果
	recommendations := append(popularBooks, latestBooks...)

	// 如果合并后的结果仍然不足，增加热门图书的数量
	if len(recommendations) < limit {
		morePopular, err := s.gorseClient.GetPopular("", limit-len(recommendations))
		if err == nil {
			recommendations = append(recommendations, morePopular...)
		}
	}

	return recommendations[:minInt(len(recommendations), limit)], nil
}

// minInt 返回两个整数中的较小值
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetPopularBooks 获取热门图书
func (s *BookService) GetPopularBooks(limit int) ([]*model.BookInfo, error) {
	// 从推荐系统获取热门图书
	popularTitles, err := s.gorseClient.GetPopular("", limit)
	if err != nil {
		return nil, fmt.Errorf("获取热门图书失败: %v", err)
	}

	// 根据标题获取完整的图书信息
	return s.getBooksByTitles(popularTitles)
}

// GetSimilarBooks 获取相似图书
func (s *BookService) GetSimilarBooks(title string, limit int) ([]*model.BookInfo, error) {
	// 从推荐系统获取相似图书
	similarTitles, err := s.gorseClient.GetItemNeighbors(title, "", limit)
	if err != nil {
		return nil, fmt.Errorf("获取相似图书失败: %v", err)
	}

	// 根据标题获取完整的图书信息
	return s.getBooksByTitles(similarTitles)
}

// getBooksByTitles 根据标题列表获取完整的图书信息
func (s *BookService) getBooksByTitles(titles []string) ([]*model.BookInfo, error) {
	if len(titles) == 0 {
		return []*model.BookInfo{}, nil
	}

	books, _, err := s.bookRepo.BatchGetBooksByTitles(titles)
	if err != nil {
		return nil, fmt.Errorf("获取图书详细信息失败: %v", err)
	}

	return books, nil
}
