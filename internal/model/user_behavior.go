package model

import "time"

// BehaviorType 用户行为类型
type BehaviorType string

const (
	BehaviorClick    BehaviorType = "click"     // 点击
	BehaviorView     BehaviorType = "view"      // 浏览
	BehaviorStayTime BehaviorType = "stay_time" // 停留时间
)

// UserBehavior 用户行为记录
type UserBehavior struct {
	ID          string       `json:"id" gorm:"primaryKey"`
	UserID      string       `json:"user_id" gorm:"not null;index"`
	BookID      string       `json:"book_id" gorm:"not null;index"`
	Type        BehaviorType `json:"type" gorm:"not null"`
	Element     string       `json:"element"`                   // 交互的元素
	Position    Position     `json:"position" gorm:"type:json"` // 点击位置
	ScrollDepth int          `json:"scroll_depth"`              // 滚动深度(%)
	StayTime    int          `json:"stay_time"`                 // 停留时间(秒)
	Timestamp   time.Time    `json:"timestamp" gorm:"not null;index"`
	Extra       string       `json:"extra"` // 额外信息（JSON格式）
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (UserBehavior) TableName() string {
	return "user_behaviors"
}

// Position 点击位置
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
