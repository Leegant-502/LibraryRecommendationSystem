package service

import (
	"fmt"
	"library/internal/gorse"
	"library/internal/model"
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

// RecordUserBehavior 统一的用户行为记录接口
func (s *BookService) RecordUserBehavior(req *UserBehaviorRequest) error {
	// 验证请求参数
	if err := req.Validate(); err != nil {
		return fmt.Errorf("参数验证失败: %w", err)
	}

	var extra map[string]interface{}
	feedbackType := req.BehaviorType

	// 根据行为类型处理额外数据
	switch req.BehaviorType {
	case "read":
		extra = map[string]interface{}{
			"read_time": *req.ReadTimeMinutes,
		}
	case "stay_time":
		extra = map[string]interface{}{
			"stay_time": *req.StayTimeSeconds,
		}
		// 根据停留时间决定反馈类型
		if *req.StayTimeSeconds >= 30 {
			feedbackType = "read"
		} else {
			feedbackType = "view"
		}
	case "view", "click":
		// 这些行为类型不需要额外处理
		extra = req.Extra
	}

	// 合并用户提供的额外信息
	if req.Extra != nil {
		if extra == nil {
			extra = make(map[string]interface{})
		}
		for k, v := range req.Extra {
			extra[k] = v
		}
	}

	// 记录到Gorse推荐系统
	return s.gorseClient.InsertFeedback(feedbackType, req.UserID, req.BookTitle, time.Now().Unix(), extra)
}

// 保留原有方法以兼容现有代码
func (s *BookService) RecordBookView(userID, title string) error {
	return s.RecordUserBehavior(&UserBehaviorRequest{
		UserID:       userID,
		BookTitle:    title,
		BehaviorType: "view",
	})
}

func (s *BookService) RecordBookClick(userID, title string) error {
	return s.RecordUserBehavior(&UserBehaviorRequest{
		UserID:       userID,
		BookTitle:    title,
		BehaviorType: "click",
	})
}

func (s *BookService) RecordBookRead(userID, title string, readTimeMinutes int) error {
	return s.RecordUserBehavior(&UserBehaviorRequest{
		UserID:          userID,
		BookTitle:       title,
		BehaviorType:    "read",
		ReadTimeMinutes: &readTimeMinutes,
	})
}

func (s *BookService) RecordBookStayTime(userID, title string, stayTimeSeconds int) error {
	return s.RecordUserBehavior(&UserBehaviorRequest{
		UserID:          userID,
		BookTitle:       title,
		BehaviorType:    "stay_time",
		StayTimeSeconds: &stayTimeSeconds,
	})
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
