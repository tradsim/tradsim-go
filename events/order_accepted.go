package events

import (
	"fmt"
	"time"

	"github.com/tradsim/tradsim-go/models"
)

// OrderAccepted defines a order accepted event
type OrderAccepted struct {
	OrderEvent
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Quantity  uint    `json:"quantity"`
	Direction string  `json:"direction"`
}

func (e *OrderAccepted) String() string {
	return fmt.Sprintf("%s %s@%f %s %d", e.OrderEvent.String(), e.Symbol, e.Price, e.Direction, e.Quantity)
}

// NewOrderAccepted creates a new order accepted event
func NewOrderAccepted(orderID string, occured time.Time, symbol string, price float64, quantity uint, direction models.TradeDirection, version uint) *OrderAccepted {

	return &OrderAccepted{*NewOrderEvent(OrderAcceptedType, orderID, occured, version), symbol, price, quantity, direction.String()}
}
