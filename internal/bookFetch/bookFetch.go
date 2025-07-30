package bookFetch

import (
	"bytes"
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

type getTokenReq struct {
	AppId string `json:"appId"`
	Code  int    `json:"code"`
}

// TokenResponse token API 响应结构
type TokenResponse struct {
	Code    int  `json:"code"`
	Success bool `json:"success"`
	Data    struct {
		Expire int    `json:"expire"`
		Token  string `json:"token"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type getBookeReq struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

func getToken() (string, error) {
	apiURL := "https://222.204.7.196:33027/api/token/getToken"

	req := getTokenReq{
		AppId: "1494732411679d",
		Code:  914,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return fmt.Sprintln("序列化请求数据失败"), err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return fmt.Sprintln("获取token失败"), err
	}
	defer resp.Body.Close().Error()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintln("读取响应内容失败 "), err
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Sprintln("解析响应内容失败"), err
	}

	if tokenResp.Code != 200 {
		return fmt.Sprintln("获取token失败"), err
	}

	return tokenResp.Data.Token, nil
	// 读取响应内容
}

// TableName 指定表名

// FetchAndSaveBooks 从API获取图书数据并保存到数据库
func fetchAndSaveBooks(db *gorm.DB) error {
	// API endpoint URL
	apiURL := "https://222.204.7.196:33027/api/dwd_ryy_xszj_tsjbxxzlb"

	//获取token
	token, err := getToken()
	if err != nil {
		return fmt.Errorf("获取token失败: %v", err)
	}

	reqBody := getBookeReq{
		PageNum:  230,
		PageSize: 1000,
	}
	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-H3C-TOKEN", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}

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

func TimelyFetchAndSaveBooks(db *gorm.DB) {
	for {
		if err := fetchAndSaveBooks(db); err != nil {
			fmt.Println("定时任务执行失败:", err)
		}
		time.Sleep(24 * time.Hour) // 每小时执行一次
	}
}
