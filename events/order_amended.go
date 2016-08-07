package events

import "time"

// OrderAmended defines a order amended event
type OrderAmended struct {
	OrderEvent
	Quantity uint `json:"quantity"`
}

// NewOrderAmended creates a new order amed pending event
func NewOrderAmended(orderID string, quantity uint, occured time.Time, version uint) *OrderAmended {

	return &OrderAmended{*NewOrderEvent(OrderAmendedType, orderID, occured, version), quantity}
}
