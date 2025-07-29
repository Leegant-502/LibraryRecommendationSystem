package repository

import (
	"library/internal/domain/model"

	"gorm.io/gorm"
)

// BookRepository 图书仓储接口
type BookRepository interface {
	CreateBookInfo(book *model.BookInfo) error
	GetBookInfoByID(id string) (*model.BookInfo, error)
	GetBookInfoByBarcode(barcode string) (*model.BookInfo, error)
	UpdateBookInfo(book *model.BookInfo) error
	FindByIDs(ids []string) ([]*model.BookInfo, error)
	BatchGetBookInfo(pageSize, pageNumber int, ids []string) ([]*model.BookInfo, int64, error)
	GetBookByTitle(title string) (*model.BookInfo, error)
	BatchGetBooksByTitles(titles []string) ([]*model.BookInfo, int64, error)
}

// PostgresBookRepository PostgreSQL实现
type PostgresBookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &PostgresBookRepository{db: db}
}

func (r *PostgresBookRepository) CreateBookInfo(book *model.BookInfo) error {
	return r.db.Create(book).Error
}

func (r *PostgresBookRepository) GetBookInfoByID(id string) (*model.BookInfo, error) {
	var book model.BookInfo
	err := r.db.Where("id = ?", id).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *PostgresBookRepository) GetBookInfoByBarcode(barcode string) (*model.BookInfo, error) {
	var book model.BookInfo
	err := r.db.Where("book_barcode = ?", barcode).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *PostgresBookRepository) UpdateBookInfo(book *model.BookInfo) error {
	return r.db.Save(book).Error
}

func (r *PostgresBookRepository) FindByIDs(ids []string) ([]*model.BookInfo, error) {
	var books []*model.BookInfo
	err := r.db.Where("id IN ?", ids).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

// BatchGetBookInfo 批量获取图书信息（支持分页）
func (r *PostgresBookRepository) BatchGetBookInfo(pageSize, pageNumber int, ids []string) ([]*model.BookInfo, int64, error) {
	var books []*model.BookInfo
	var total int64

	// 计算总数
	if err := r.db.Model(&model.BookInfo{}).Where("id IN ?", ids).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (pageNumber - 1) * pageSize
	err := r.db.Where("id IN ?", ids).
		Offset(offset).
		Limit(pageSize).
		Find(&books).Error

	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (r *PostgresBookRepository) GetBookByTitle(title string) (*model.BookInfo, error) {
	var book model.BookInfo
	err := r.db.Where("title = ?", title).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *PostgresBookRepository) BatchGetBooksByTitles(titles []string) ([]*model.BookInfo, int64, error) {
	var books []*model.BookInfo
	var total int64

	err := r.db.Model(&model.BookInfo{}).Where("title IN ?", titles).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("title IN ?", titles).Find(&books).Error
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}
