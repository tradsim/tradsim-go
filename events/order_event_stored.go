package events

import "time"

// OrderEventStored defines a order event stored
type OrderEventStored struct {
	OrderEvent
}

func (e *OrderEventStored) String() string {
	return e.OrderEvent.String()
}

// NewOrderEventStored creates a new order event stored event
func NewOrderEventStored(orderID string, occured time.Time, version uint) *OrderEventStored {

	return &OrderEventStored{*NewOrderEvent(OrderEventStoredType, orderID, occured, version)}
}
