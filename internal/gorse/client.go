package gorse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client Gorse API 客户端
type Client struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

// NewClient 创建新的 Gorse 客户端
func NewClient(endpoint, apiKey string) *Client {
	return &Client{
		endpoint: endpoint,
		apiKey:   apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// InsertFeedback 插入用户反馈数据
func (c *Client) InsertFeedback(feedbackType, userID, itemID string, timestamp int64, extra map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/feedback", c.endpoint)

	feedback := map[string]interface{}{
		"FeedbackType": feedbackType,
		"UserId":       userID,
		"ItemId":       itemID,
		"Timestamp":    timestamp,
	}

	if extra != nil {
		feedback["Extra"] = extra
	}

	jsonData, err := json.Marshal(feedback)
	if err != nil {
		return fmt.Errorf("序列化反馈数据失败: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API返回错误状态码: %d", resp.StatusCode)
	}

	return nil
}

// GetRecommend 获取个性化推荐
func (c *Client) GetRecommend(userID string, category string, n int) ([]string, error) {
	url := fmt.Sprintf("%s/api/recommend/%s?n=%d", c.endpoint, userID, n)
	if category != "" {
		url += "&category=" + category
	}
	return c.getItems(url)
}

// GetPopular 获取热门图书
func (c *Client) GetPopular(category string, n int) ([]string, error) {
	url := fmt.Sprintf("%s/api/popular?n=%d", c.endpoint, n)
	if category != "" {
		url += "&category=" + category
	}
	return c.getItems(url)
}

// GetItemNeighbors 获取相似图书
func (c *Client) GetItemNeighbors(itemID string, category string, n int) ([]string, error) {
	url := fmt.Sprintf("%s/api/item/%s/neighbors?n=%d", c.endpoint, itemID, n)
	if category != "" {
		url += "&category=" + category
	}
	return c.getItems(url)
}

// GetLatest 获取最新图书
func (c *Client) GetLatest(category string, n int) ([]string, error) {
	url := fmt.Sprintf("%s/api/latest?n=%d", c.endpoint, n)
	if category != "" {
		url += "&category=" + category
	}
	return c.getItems(url)
}

// getItems 通用的获取项目列表方法
func (c *Client) getItems(url string) ([]string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误状态码: %d", resp.StatusCode)
	}

	var items []string
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

// Feedback 用户反馈结构
type Feedback struct {
	FeedbackType string    `json:"FeedbackType"`
	UserId       string    `json:"UserId"`
	ItemId       string    `json:"ItemId"`
	Timestamp    time.Time `json:"Timestamp"`
}
