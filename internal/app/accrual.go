package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rawen554/go-loyal/internal/models"
)

const OrdersAPI = "/api/orders/{number}"

var ErrNoOrder = errors.New("order is not processed")
var ErrServiceBisy = errors.New("accrual is bisy")

type AccrualClient struct {
	client *resty.Client
}

type Accrual interface {
	GetOrderInfo(num string) (*AccrualOrderInfoShema, error)
}

type AccrualOrderInfoShema struct {
	Order   string        `json:"order"`
	Status  models.Status `json:"status"`
	Accrual float64       `json:"accrual,omitempty"`
}

func NewAccrualClient(accrualAddr string) (Accrual, error) {
	return &AccrualClient{
		client: resty.New().SetBaseURL(accrualAddr),
	}, nil
}

func (a *AccrualClient) GetOrderInfo(num string) (*AccrualOrderInfoShema, error) {
	var orderInfo AccrualOrderInfoShema
	result, err := a.client.R().SetResult(&orderInfo).SetPathParam("number", num).Get(OrdersAPI)
	if err != nil {
		return nil, fmt.Errorf("error getting order info from accrual: %w", err)
	}

	switch result.StatusCode() {
	case http.StatusOK:
		return &orderInfo, nil
	case http.StatusNoContent:
		return nil, ErrNoOrder
	case http.StatusTooManyRequests:
		return nil, ErrServiceBisy
	default:
		return nil, fmt.Errorf("unknown exception: %w", err)
	}
}
