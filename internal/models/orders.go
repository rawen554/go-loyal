package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var ErrOrderHasBeenProcessedByUser = errors.New("this order already been processed by user")
var ErrUserHasNoOrders = errors.New("this user has no orders")

type Status string

const (
	NEW        Status = "NEW"
	PROCESSING Status = "PROCESSING"
	INVALID    Status = "INVALID"
	PROCESSED  Status = "PROCESSED"
)

func (s *Status) Scan(value interface{}) error {
	sv, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal Status value: ", value))
	}

	*s = Status(sv)
	return nil
	// *s = Status(value.([]byte))
	// return nil
}

func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

type OrderTime time.Time

func (ot OrderTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(ot).Format(time.RFC3339))
	return []byte(formatted), nil
}

type Order struct {
	Number     uint64    `gorm:"primaryKey" json:"number"`
	UploadedAt OrderTime `gorm:"default:now()" json:"uploaded_at"`
	UserID     uint64    `json:"-"`
	Accrual    uint64    `json:"accrual,omitempty"`
	User       User      `json:"-"`
	Status     Status    `sql:"type:order_status" json:"status"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.Status = NEW
	return nil
}
