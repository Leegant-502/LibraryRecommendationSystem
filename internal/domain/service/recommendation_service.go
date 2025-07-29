package service

import (
	"library/internal/domain/model"
	"time"
)

// RecommendationService 推荐领域服务接口
type RecommendationService interface {
	// 用户行为相关
	RecordBookView(userID, bookID string) error
	GetUserBehaviors(userID string, startTime, endTime time.Time) ([]*model.UserBehavior, error)

	// 推荐相关
	GetPersonalizedRecommendations(userID string, limit int) (*model.RecommendationResponse, error)
	GetCategoryRecommendations(limit int) ([]model.RecommendationCategory, error)
	GetSimilarBooks(bookID string, limit int) (*model.RecommendationResponse, error)
}

// BookService 图书领域服务接口
type BookService interface {
	// 图书基本操作
	CreateBook(book *model.BookInfo) error
	GetBook(id string) (*model.BookInfo, error)
	UpdateBook(book *model.BookInfo) error
	DeleteBook(id string) error

	// 批量操作
	BatchGetBooks(req *model.BatchBookRequest) (*model.BatchBookResponse, error)
	BatchCreateBooks(books []*model.BookInfo) error

	// 外部数据同步
	FetchAndSaveBookInfo(id string) error
	BatchFetchAndSaveBookInfo(ids []string) error
}
