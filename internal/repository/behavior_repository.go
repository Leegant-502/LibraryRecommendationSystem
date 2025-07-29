package repository

import (
	"gorm.io/gorm"
	"library/internal/domain/model"
	"time"
)

// BehaviorRepository 用户行为仓储接口
type BehaviorRepository interface {
	Save(behavior *model.UserBehavior) error
	FindByUserIDAndTimeRange(userID string, startTime, endTime time.Time) ([]*model.UserBehavior, error)
	FindByBookIDAndTimeRange(bookID string, startTime, endTime time.Time) ([]*model.UserBehavior, error)
	FindClicksByBookIDAndTimeRange(bookID string, startTime, endTime time.Time) ([]*model.UserBehavior, error)
}

// PostgresBehaviorRepository PostgreSQL实现
type PostgresBehaviorRepository struct {
	db *gorm.DB
}

func NewPostgresBehaviorRepository(db *gorm.DB) BehaviorRepository {
	return &PostgresBehaviorRepository{db: db}
}

func (r *PostgresBehaviorRepository) Save(behavior *model.UserBehavior) error {
	return r.db.Create(behavior).Error
}

func (r *PostgresBehaviorRepository) FindByUserIDAndTimeRange(userID string, startTime, endTime time.Time) ([]*model.UserBehavior, error) {
	var behaviors []*model.UserBehavior
	err := r.db.Where("user_id = ? AND timestamp BETWEEN ? AND ?", userID, startTime, endTime).
		Find(&behaviors).Error
	return behaviors, err
}

func (r *PostgresBehaviorRepository) FindByBookIDAndTimeRange(bookID string, startTime, endTime time.Time) ([]*model.UserBehavior, error) {
	var behaviors []*model.UserBehavior
	err := r.db.Where("book_id = ? AND timestamp BETWEEN ? AND ?", bookID, startTime, endTime).
		Find(&behaviors).Error
	return behaviors, err
}

func (r *PostgresBehaviorRepository) FindClicksByBookIDAndTimeRange(bookID string, startTime, endTime time.Time) ([]*model.UserBehavior, error) {
	var behaviors []*model.UserBehavior
	err := r.db.Where("book_id = ? AND type = ? AND timestamp BETWEEN ? AND ?",
		bookID, model.BehaviorClick, startTime, endTime).
		Find(&behaviors).Error
	return behaviors, err
}
