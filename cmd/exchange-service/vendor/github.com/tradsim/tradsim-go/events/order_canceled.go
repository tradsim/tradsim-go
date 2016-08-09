package events

import "time"

// OrderCanceled defines a order canceled event
type OrderCanceled struct {
	OrderEvent
}

// NewOrderCanceled creates a new order amed pending event
func NewOrderCanceled(orderID string, occured time.Time, version uint) *OrderCanceled {

	return &OrderCanceled{*NewOrderEvent(OrderCanceledType, orderID, occured, version)}
}
