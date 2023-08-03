package models

type Withdraw struct {
	OrderNum    string    `json:"order"`
	UserID      uint64    `gorm:"column:user_id" json:"-"`
	User        User      `json:"-"`
	Sum         float64   `json:"sum"`
	ProcessedAt OrderTime `gorm:"default:now()" json:"processed_at"`
}

func (w *Withdraw) TableName() string {
	return "withdrawals"
}

type BalanceWithdrawShema struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
