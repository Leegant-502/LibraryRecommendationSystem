package bookFetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// BookInfo 图书信息结构体
type BookInfo struct {
	ID                   string    `json:"id" gorm:"primaryKey"`
	BookID               string    `json:"tsbh"`  // 图书编号
	BookBarcode          string    `json:"tstxm"` // 图书条形码
	Title                string    `json:"zbt"`   // 正标题
	PublicationNumber    string    `json:"tscbh"` // 图书出版号
	PrimaryAuthor        string    `json:"dyzz"`  // 第一作者
	ClassificationNumber string    `json:"flh"`   // 分类号
	LanguageCode         string    `json:"yzm"`   // 语种码
	Edition              string    `json:"bc"`    // 版次
	Publisher            string    `json:"cbs"`   // 出版社
	PublicationPlace     string    `json:"cbd"`   // 出版地
	PublicationDate      string    `json:"cbrq"`  // 出版日期
	DistributionUnit     string    `json:"fxdw"`  // 发行单位
	Notes                string    `json:"bz"`    // 备注
	CreatedAt            time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (BookInfo) TableName() string {
	return "book_information"
}

// FetchAndSaveBooks 从API获取图书数据并保存到数据库
func FetchAndSaveBooks(db *gorm.DB) error {
	// API endpoint URL
	apiURL := "https://222.204.7.196:33027/api/dwd_ryy_xszj_tsjbxxzlb"

	// 发送HTTP GET请求
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("获取图书数据失败: %v", err)
	}
	defer resp.Body.Close().Error()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 解析JSON数据到图书结构体切片
	var books []BookInfo
	if err := json.Unmarshal(body, &books); err != nil {
		return fmt.Errorf("解析JSON数据失败: %v", err)
	}

	// 批量保存图书数据到数据库
	if err := db.CreateInBatches(books, 1000).Error; err != nil {
		return fmt.Errorf("批量保存图书数据失败: %v", err)
	}

	fmt.Printf("成功保存 %d 本图书到数据库\n", len(books))
	return nil
}
