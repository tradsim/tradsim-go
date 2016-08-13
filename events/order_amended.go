package events

import (
	"fmt"
	"time"
)

// OrderAmended defines a order amended event
type OrderAmended struct {
	OrderEvent
	Quantity uint `json:"quantity"`
}

func (e *OrderAmended) String() string {
	return fmt.Sprintf("%s %d", e.OrderEvent.String(), e.Quantity)
}

// NewOrderAmended creates a new order amed pending event
func NewOrderAmended(orderID string, quantity uint, occured time.Time, version uint) *OrderAmended {

	return &OrderAmended{*NewOrderEvent(OrderAmendedType, orderID, occured, version), quantity}
}
