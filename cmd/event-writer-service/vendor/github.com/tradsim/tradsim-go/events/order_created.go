package events

import (
	"fmt"
	"time"
	"github.com/tradsim/tradsim-go/models"
)

// OrderCreated defines a order created event
type OrderCreated struct {
	OrderEvent
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Quantity  uint    `json:"quantity"`
	Direction string  `json:"direction"`
}

func (e *OrderCreated) String() string {
	return fmt.Sprintf("%s %s@%f %s %d", e.OrderEvent.String(), e.Symbol, e.Price, e.Direction, e.Quantity)
}

// NewOrderCreated creates a new order created event
func NewOrderCreated(orderID string, occured time.Time, symbol string, price float64, quantity uint, direction models.TradeDirection, version uint) *OrderCreated {

	return &OrderCreated{*NewOrderEvent(OrderCreatedType, orderID, occured, version), symbol, price, quantity, direction.String()}
}
