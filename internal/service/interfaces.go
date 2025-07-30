package service

import (
	"fmt"
	"library/internal/model"
)

// UserBehaviorRequest 用户行为请求
type UserBehaviorRequest struct {
	UserID          string                 `json:"user_id"`
	BookTitle       string                 `json:"book_title"`
	BehaviorType    string                 `json:"behavior_type"`
	StayTimeSeconds *int                   `json:"stay_time_seconds,omitempty"`
	ReadTimeMinutes *int                   `json:"read_time_minutes,omitempty"`
	Extra           map[string]interface{} `json:"extra,omitempty"`
}

// Validate 验证请求参数
func (r *UserBehaviorRequest) Validate() error {
	if r.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if r.BookTitle == "" {
		return fmt.Errorf("book_title is required")
	}
	if r.BehaviorType == "" {
		return fmt.Errorf("behavior_type is required")
	}

	// 验证行为类型特定的参数
	switch r.BehaviorType {
	case "read":
		if r.ReadTimeMinutes == nil {
			return fmt.Errorf("read_time_minutes is required for read behavior")
		}
	case "stay_time":
		if r.StayTimeSeconds == nil {
			return fmt.Errorf("stay_time_seconds is required for stay_time behavior")
		}
	case "view", "click":
		// 这些行为类型不需要额外参数
	default:
		return fmt.Errorf("unsupported behavior type: %s", r.BehaviorType)
	}

	return nil
}

// BookServiceInterface 图书服务接口
// 定义了图书推荐系统的核心业务能力
type BookServiceInterface interface {
	// 统一的用户行为记录接口
	RecordUserBehavior(req *UserBehaviorRequest) error

	// 推荐获取
	GetRecommendations(userID string, limit int) ([]*model.BookInfo, error)
	GetPopularBooks(limit int) ([]*model.BookInfo, error)
	GetSimilarBooks(title string, limit int) ([]*model.BookInfo, error)
}

// RecommendationServiceInterface 推荐服务接口
// 专门处理推荐算法相关的业务逻辑
type RecommendationServiceInterface interface {
	// 获取个性化推荐
	GetPersonalizedRecommendations(userID string, limit int) ([]*model.BookInfo, error)

	// 获取热门推荐
	GetPopularRecommendations(limit int) ([]*model.BookInfo, error)

	// 获取相似物品推荐
	GetSimilarItemRecommendations(itemID string, limit int) ([]*model.BookInfo, error)

	// 记录用户反馈
	RecordUserFeedback(userID, itemID, feedbackType string, timestamp int64, extra map[string]interface{}) error
}

// BehaviorTrackingServiceInterface 行为追踪服务接口
// 专门处理用户行为数据的收集和分析
type BehaviorTrackingServiceInterface interface {
	// 记录用户行为
	TrackUserBehavior(userID, itemID, behaviorType string, metadata map[string]interface{}) error

	// 批量记录用户行为
	BatchTrackUserBehavior(behaviors []model.UserBehavior) error

	// 获取用户行为历史
	GetUserBehaviorHistory(userID string, limit int) ([]*model.UserBehavior, error)

	// 分析用户兴趣
	AnalyzeUserInterests(userID string) (map[string]float64, error)
}
