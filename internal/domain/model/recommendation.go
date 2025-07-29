package model

// RecommendationResponse 推荐响应结构
type RecommendationResponse struct {
	Books    []*BookInfo `json:"books"`    // 推荐的图书列表
	Reason   string      `json:"reason"`   // 推荐原因
	Score    float64     `json:"score"`    // 推荐分数
	Category string      `json:"category"` // 推荐类别（如"热门"、"相似"、"个性化"等）
}

// RecommendationCategory 推荐类别
type RecommendationCategory struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Books       []*BookInfo `json:"books"`
	Total       int         `json:"total"`
}
