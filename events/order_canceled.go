package events

import "time"

// OrderCancelled defines a order cancelled event
type OrderCancelled struct {
	OrderEvent
}

func (e *OrderCancelled) String() string {
	return e.OrderEvent.String()
}

// NewOrderCancelled creates a new order amed pending event
func NewOrderCancelled(orderID string, occured time.Time, version uint) *OrderCancelled {

	return &OrderCancelled{*NewOrderEvent(OrderCancelledType, orderID, occured, version)}
}
