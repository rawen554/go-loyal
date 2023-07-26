package accrual

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rawen554/go-loyal/internal/models"
)

const OrdersAPI = "/api/orders/{number}"

var ErrNoOrder = errors.New("order is not processed")
var ErrServiceBisy = errors.New("accrual is bisy")

type ServiceBisyError struct {
	CoolDown time.Duration
	MaxRPM   int
	Err      error
}

func (sbe *ServiceBisyError) Error() string {
	return fmt.Sprintf("wait: %vs; max rpm: %v; %v", sbe.CoolDown.Seconds(), sbe.MaxRPM, sbe.Err)
}

func NewServiceBisyError(cooldown time.Duration, rpm int, err error) error {
	return &ServiceBisyError{
		CoolDown: cooldown,
		MaxRPM:   rpm,
		Err:      err,
	}
}

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

// retryable (hashicorp) 5xx

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
		cooldown, err := strconv.Atoi(result.Header().Get("Retry-After"))
		if err != nil {
			log.Printf("error converting header Retry-After: %v", err)
			return nil, err
		}

		r := regexp.MustCompile(`(\d+)`)
		rpm := r.Find(result.Body())
		if rpm == nil {
			log.Printf("not found MaxRPM in body: %v", err)
			return nil, err
		}
		preparedRPM, err := strconv.Atoi(string(rpm))
		if err != nil {
			log.Printf("cant convert MaxRPM to int: %v", err)
			return nil, err
		}

		return nil, NewServiceBisyError(time.Duration(cooldown)*time.Second, preparedRPM, err)
	default:
		return nil, fmt.Errorf("unknown exception: %w", err)
	}
}