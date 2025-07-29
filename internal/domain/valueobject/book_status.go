package valueobject

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

// IsAvailable 检查图书是否可借阅
func (s BookStatus) IsAvailable() bool {
	return s == BookStatusAvailable
}

// IsBorrowed 检查图书是否已借出
func (s BookStatus) IsBorrowed() bool {
	return s == BookStatusBorrowed
}
