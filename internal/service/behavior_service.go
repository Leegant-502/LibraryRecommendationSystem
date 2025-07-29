package service

import (
	"library/internal/domain/model"
	"library/internal/gorse"
	"library/internal/repository"
	"time"
)

// BehaviorService 用户行为服务
type BehaviorService struct {
	behaviorRepo repository.BehaviorRepository
	gorseClient  *gorse.Client
}

// NewBehaviorService 创建行为服务实例
func NewBehaviorService(behaviorRepo repository.BehaviorRepository, gorseClient *gorse.Client) *BehaviorService {
	return &BehaviorService{
		behaviorRepo: behaviorRepo,
		gorseClient:  gorseClient,
	}
}

// TrackBehavior 记录用户行为
func (s *BehaviorService) TrackBehavior(behavior *model.UserBehavior) error {
	// 1. 保存行为数据
	if err := s.behaviorRepo.Save(behavior); err != nil {
		return err
	}

	// 2. 根据不同行为类型处理
	switch behavior.Type {
	case model.BehaviorView:
		// 记录浏览行为到推荐系统
		return s.gorseClient.InsertFeedback("view", behavior.UserID, behavior.BookID, behavior.Timestamp.Unix(), nil)

	case model.BehaviorStayTime:
		if behavior.StayTime >= 30 { // 停留超过30秒视为深度阅读
			extra := map[string]interface{}{
				"stay_time": behavior.StayTime,
			}
			return s.gorseClient.InsertFeedback("read", behavior.UserID, behavior.BookID, behavior.Timestamp.Unix(), extra)
		}
	}

	return nil
}

// GetUserBehaviors 获取用户行为历史
func (s *BehaviorService) GetUserBehaviors(userID string, startTime, endTime time.Time) ([]*model.UserBehavior, error) {
	return s.behaviorRepo.FindByUserIDAndTimeRange(userID, startTime, endTime)
}

// GetBookBehaviors 获取图书行为统计
func (s *BehaviorService) GetBookBehaviors(bookID string, startTime, endTime time.Time) (map[model.BehaviorType]int, error) {
	behaviors, err := s.behaviorRepo.FindByBookIDAndTimeRange(bookID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	stats := make(map[model.BehaviorType]int)
	for _, b := range behaviors {
		stats[b.Type]++
	}
	return stats, nil
}
