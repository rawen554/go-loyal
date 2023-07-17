package models

type Withdraw struct {
	OrderNum    string
	Order       Order `gorm:"foreignKey:OrderNum" json:"-"`
	UserID      uint64
	User        User `json:"-"`
	Sum         float64
	ProcessedAt OrderTime
}
