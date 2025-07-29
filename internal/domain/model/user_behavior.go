package model

import "time"

// BehaviorType 用户行为类型
type BehaviorType string

const (
	BehaviorClick    BehaviorType = "click"     // 点击
	BehaviorView     BehaviorType = "view"      // 浏览
	BehaviorStayTime BehaviorType = "stay_time" // 停留时间
	BehaviorScroll   BehaviorType = "scroll"    // 滚动
)

// UserBehavior 用户行为记录
type UserBehavior struct {
	ID          string       `json:"id"`
	UserID      string       `json:"user_id"`
	BookID      string       `json:"book_id"`
	Type        BehaviorType `json:"type"`
	Element     string       `json:"element"`      // 交互的元素
	Position    Position     `json:"position"`     // 点击位置
	ScrollDepth int          `json:"scroll_depth"` // 滚动深度(%)
	StayTime    int          `json:"stay_time"`    // 停留时间(秒)
	Timestamp   time.Time    `json:"timestamp"`
	Extra       string       `json:"extra"` // 额外信息（JSON格式）
}

// Position 点击位置
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
