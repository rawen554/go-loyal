package processing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rawen554/go-loyal/internal/adapters/accrual"
	"github.com/rawen554/go-loyal/internal/adapters/store"
	"github.com/rawen554/go-loyal/internal/models"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type ProcessingController struct {
	ordersChan chan *models.Order
	store      store.Store
	accrual    accrual.Accrual
	errg       *errgroup.Group
	logger     *zap.SugaredLogger
}

const chanLen = 10

func NewProcessingController(store store.Store, accrual accrual.Accrual) *ProcessingController {
	ordersChan := make(chan *models.Order, chanLen)

	instance := &ProcessingController{
		ordersChan: ordersChan,
		store:      store,
		accrual:    accrual,
		errg:       new(errgroup.Group),
	}

	go instance.ListenOrders()

	return instance
}

func (p *ProcessingController) ListenOrders() {
	ticker := time.NewTicker(chanLen * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		orders, err := p.store.GetUnprocessedOrders()
		if err != nil {
			p.logger.Errorf("error getting unprocessed orders from store: %v", err)
		}
		for i := range orders {
			p.ordersChan <- &orders[i]
		}
	}
}

func (p *ProcessingController) Process(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case o := <-p.ordersChan:
				p.errg.Go(func() error {
					if o.Status == models.NEW {
						_, err := p.store.UpdateOrder(&models.Order{Number: o.Number, Status: models.PROCESSING})
						if err != nil {
							return fmt.Errorf("error updating order from accrual: %w", err)
						}
					}

					info, err := p.accrual.GetOrderInfo(o.Number)
					if err != nil {
						return fmt.Errorf("error getting order: %w", err)
					}

					if info.Status == models.PROCESSED || info.Status == models.INVALID {
						_, err := p.store.UpdateOrder(
							&models.Order{
								Number:  info.Order,
								UserID:  o.UserID,
								Accrual: info.Accrual,
								Status:  info.Status,
							})
						if err != nil {
							return fmt.Errorf("error updating order: %w", err)
						}
					}

					return nil
				})
			}
		}
	}(ctx)

	if err := p.errg.Wait(); err != nil {
		var serviceBisyError *accrual.ServiceBisyError
		if errors.As(err, &serviceBisyError) {
			p.logger.Infof("service bisy: %v", serviceBisyError)
			time.Sleep(serviceBisyError.CoolDown)
		} else {
			p.logger.Errorf("unhandled error: %v", err)
		}
	}
}
