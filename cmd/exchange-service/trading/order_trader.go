package trading

import (
	"github.com/mantzas/adaptlog"
	"github.com/tradsim/tradsim-go/models"
)

// Trader interface
type Trader interface {
	Trade(book *models.OrderBook, order *models.Order) (*models.Order, error)
}

// OrderTrader implementation
type OrderTrader struct {
	logger adaptlog.LevelLogger
}

// Trade processes a order against the book
func (ot *OrderTrader) Trade(book *models.OrderBook, order *models.Order) (*models.Order, error) {

	return order, nil
}
