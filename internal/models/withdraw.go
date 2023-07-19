package models

import "gorm.io/gorm"

type Withdraw struct {
	OrderNum    string    `json:"order"`
	Order       Order     `gorm:"foreignKey:OrderNum" json:"-"`
	UserID      uint64    `json:"-"`
	User        User      `json:"-"`
	Sum         float64   `json:"sum"`
	ProcessedAt OrderTime `gorm:"default:now()" json:"processed_at"`
}

type BalanceWithdrawShema struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (w *Withdraw) AfterCreate(tx *gorm.DB) (err error) {
	result := tx.Model(&User{}).Where("id = ?", w.UserID).
		Updates(map[string]interface{}{"balance": gorm.Expr("balance - ?", w.Sum), "withdrawn": gorm.Expr("withdrawn + ?", w.Sum)})
	return result.Error
}
