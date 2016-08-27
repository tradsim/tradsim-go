package trading

import (
	"sync"

	"github.com/mantzas/adaptlog"
	"github.com/tradsim/tradsim-go/models"
)

// Amender interface
type Amender interface {
	Amend(book *models.OrderBook, order *models.Order) bool
}

// OrderAmender amends a order in the book
type OrderAmender struct {
	mu     sync.Mutex
	logger adaptlog.LevelLogger
}

// NewOrderAmender creates a new order amender
func NewOrderAmender() *OrderAmender {
	return &OrderAmender{sync.Mutex{}, adaptlog.NewStdLevelLogger("OrderAmender")}
}

// Amend a order in the order book
func (oa *OrderAmender) Amend(book *models.OrderBook, order *models.Order) bool {

	oa.mu.Lock()
	defer oa.mu.Unlock()

	prices, ok := book.Symbols[order.Symbol]

	if !ok {
		oa.logger.Errorf("Symbol %s not found", order.Symbol)
		return false
	}

	found, i := findPrice(prices, order.Price)

	if !found {
		oa.logger.Errorf("Price %f not found", order.Price)
		return false
	}

	if order.Direction == models.Buy {

		for _, o := range prices[i].Buy.Orders {

			if o.ID == order.ID {
				prices[i].Buy.Quantity += order.Quantity
				o.Amend(order.Quantity)
				return true
			}
		}

	} else {

		for _, o := range prices[i].Sell.Orders {

			if o.ID == order.ID {
				prices[i].Sell.Quantity += order.Quantity
				o.Amend(order.Quantity)
				return true
			}
		}
	}

	return false
}
