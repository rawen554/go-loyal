package app

import (
	"github.com/go-resty/resty/v2"
	"github.com/rawen554/go-loyal/internal/models"
)

const OrdersApi = "/api/orders/{number}"

type AccrualClient struct {
	client *resty.Client
}

type AccrualOrderInfoShema struct {
	Order   string        `json:"order"`
	Status  models.Status `json:"status"`
	Accrual uint64        `json:"accrual,omitempty"`
}

func NewAccrualClient(accrualAddr string) (*AccrualClient, error) {
	return &AccrualClient{
		client: resty.New().SetBaseURL(accrualAddr),
	}, nil
}

func (a *AccrualClient) GetOrderInfo(num string) (*AccrualOrderInfoShema, error) {
	var orderInfo AccrualOrderInfoShema
	_, err := a.client.R().SetResult(&orderInfo).SetPathParam("number", num).Get(OrdersApi)
	if err != nil {
		return nil, err
	}

	return &orderInfo, nil
}
