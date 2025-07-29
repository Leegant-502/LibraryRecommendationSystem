package model

// BatchBookRequest 批量获取图书请求
type BatchBookRequest struct {
	BookIDs    []string `json:"book_ids"`
	PageSize   int      `json:"page_size"`
	PageNumber int      `json:"page_number"`
}

// BatchBookResponse 批量获取图书响应
type BatchBookResponse struct {
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Books    []*BookInfo `json:"books"` // 使用指针切片
	HasMore  bool        `json:"has_more"`
}
