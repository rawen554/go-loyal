package processing

import (
	"context"
	"log"
	"time"

	"github.com/rawen554/go-loyal/internal/adapters/accrual"
	"github.com/rawen554/go-loyal/internal/adapters/store"
	"github.com/rawen554/go-loyal/internal/models"
)

type ProcessingController struct {
	ordersChan chan *models.Order
	store      store.Store
	accrual    accrual.Accrual
}

func NewProcessingController(ordersChan chan *models.Order, store store.Store, accrual accrual.Accrual) *ProcessingController {
	instance := &ProcessingController{
		ordersChan: ordersChan,
		store:      store,
		accrual:    accrual,
	}

	go instance.ListenOrders()

	return instance
}

func (p *ProcessingController) ListenOrders() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		orders, err := p.store.GetUnprocessedOrders()
		if err != nil {
			log.Printf("error getting unprocessed orders from store: %v", err)
		}
		for i, _ := range orders {
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
				if o.Status == models.NEW {
					go p.store.UpdateOrder(&models.Order{Number: o.Number, Status: models.PROCESSING})
				}

				info, err := p.accrual.GetOrderInfo(o.Number)
				if err != nil {
					log.Printf("error interacting with accrual: %v", err)
					return
				}

				if info.Status == models.PROCESSED || info.Status == models.INVALID {
					go func() {
						rowsAffected, err := p.store.UpdateOrder(&models.Order{Number: info.Order, UserID: o.UserID, Accrual: info.Accrual, Status: info.Status})
						if rowsAffected == 0 || err != nil {
							log.Printf("error updating order: %v", err)
						}
					}()
				}
			}
		}
	}(ctx)
}
