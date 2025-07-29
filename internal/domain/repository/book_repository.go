package repository

import (
	"library/internal/domain/model"
	"time"
)

// BookRepository 图书仓储接口
type BookRepository interface {
	CreateBookInfo(book *model.BookInfo) error
	GetBookInfoByID(id string) (*model.BookInfo, error)
	GetBookInfoByBarcode(barcode string) (*model.BookInfo, error)
	UpdateBookInfo(book *model.BookInfo) error
	FindByIDs(ids []string) ([]*model.BookInfo, error)
	BatchGetBookInfo(pageSize, pageNumber int, ids []string) ([]*model.BookInfo, int64, error)
}

// RecommendationRepository 推荐仓储接口
type RecommendationRepository interface {
	SaveUserBehavior(behavior *model.UserBehavior) error
	GetUserBehaviors(userID string, startTime, endTime time.Time) ([]*model.UserBehavior, error)
	GetPopularBooks(limit int) ([]*model.BookInfo, error)
	GetLatestBooks(limit int) ([]*model.BookInfo, error)
}
