package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// BookStatus 图书状态枚举
type BookStatus string

const (
	// BookStatusAvailable 可借阅
	BookStatusAvailable BookStatus = "available"
	// BookStatusBorrowed 已借出
	BookStatusBorrowed BookStatus = "borrowed"
	// BookStatusReserved 已预约
	BookStatusReserved BookStatus = "reserved"
	// BookStatusMaintenance 维护中
	BookStatusMaintenance BookStatus = "maintenance"
	// BookStatusLost 丢失
	BookStatusLost BookStatus = "lost"
	// BookStatusDamaged 损坏
	BookStatusDamaged BookStatus = "damaged"
)

// String 返回状态的字符串表示
func (s BookStatus) String() string {
	return string(s)
}

// IsValid 检查状态是否有效
func (s BookStatus) IsValid() bool {
	switch s {
	case BookStatusAvailable, BookStatusBorrowed, BookStatusReserved,
		BookStatusMaintenance, BookStatusLost, BookStatusDamaged:
		return true
	default:
		return false
	}
}

// 最外层响应结构体
type Response struct {
	Msg     string   `json:"msg"`     // 消息提示
	Code    string   `json:"code"`    // 状态码
	Data    Data     `json:"data"`    // 分页数据
	Records []Record `json:"records"` // 记录列表
}

// 分页信息结构体
type Data struct {
	Current int `json:"current"` // 当前页码
	Size    int `json:"size"`    // 每页条数
	Total   int `json:"total"`   // 总记录数
	Pages   int `json:"pages"`   // 总页数
}

// Record 记录结构体
type Record struct {
	ID     int         `json:"id"`     // 记录ID
	Type   string      `json:"type"`   // 记录类型
	Detail interface{} `json:"detail"` // 记录详情
}

// BookInfo 图书信息表
type BookInfo struct {
	ID                   decimal.Decimal `json:"id" gorm:"primaryKey"`
	BookID               string          `json:"book_id" gorm:"unique;column:book_id"`
	BookBarcode          string          `json:"book_barcode" gorm:"unique"`
	Title                string          `json:"title" gorm:"not null"`
	PublicationNumber    string          `json:"publication_number"`
	PrimaryAuthor        string          `json:"primary_author"`
	ClassificationNumber string          `json:"classification_number"`
	LanguageCode         string          `json:"language_code"`
	Edition              string          `json:"edition"`
	Publisher            string          `json:"publisher"`
	PublicationPlace     string          `json:"publication_place"`
	PublicationDate      time.Time       `json:"publication_date"`
	DistributionUnit     string          `json:"distribution_unit"`
	Notes                string          `json:"notes"`
	Status               BookStatus      `json:"status" gorm:"type:varchar(20);default:'available'"`
	CreatedAt            time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (BookInfo) TableName() string {
	return "book_information"
}

// APIBookInfo 从API获取的图书信息结构
type APIBookInfo struct {
	ID                   string `json:"id"`
	BookID               string `json:"tsbh"`  // 图书编号
	BookBarcode          string `json:"tstxm"` // 图书条形码
	Title                string `json:"zbt"`   // 正标题
	PublicationNumber    string `json:"tscbh"` // 图书出版号
	PrimaryAuthor        string `json:"dyzz"`  // 第一作者
	ClassificationNumber string `json:"flh"`   // 分类号
	LanguageCode         string `json:"yzm"`   // 语种码
	Edition              string `json:"bc"`    // 版次
	Publisher            string `json:"cbs"`   // 出版社
	PublicationPlace     string `json:"cbd"`   // 出版地
	PublicationDate      string `json:"cbrq"`  // 出版日期
	DistributionUnit     string `json:"fxdw"`  // 发行单位
	Notes                string `json:"bz"`    // 备注
}

// ToBookInfo 将API数据转换为数据库模型
func (a *APIBookInfo) ToBookInfo() (*BookInfo, error) {
	id, err := decimal.NewFromString(a.ID)
	if err != nil {
		return nil, err
	}

	pubDate, err := time.Parse("2006-01-02", a.PublicationDate)
	if err != nil {
		// 如果日期解析失败，使用零值
		pubDate = time.Time{}
	}

	return &BookInfo{
		ID:                   id,
		BookID:               a.BookID,
		BookBarcode:          a.BookBarcode,
		Title:                a.Title,
		PublicationNumber:    a.PublicationNumber,
		PrimaryAuthor:        a.PrimaryAuthor,
		ClassificationNumber: a.ClassificationNumber,
		LanguageCode:         a.LanguageCode,
		Edition:              a.Edition,
		Publisher:            a.Publisher,
		PublicationPlace:     a.PublicationPlace,
		PublicationDate:      pubDate,
		DistributionUnit:     a.DistributionUnit,
		Notes:                a.Notes,
	}, nil
}
