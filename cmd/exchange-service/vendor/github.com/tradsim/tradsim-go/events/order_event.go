package events

import (
	"fmt"
	"time"
)

// OrderEvent defines a order base event
type OrderEvent struct {
	EventType OrderEventType `json:"event_type"`
	OrderID   string         `json:"id"`
	Occured   time.Time      `json:"occured"`
	Version   uint           `json:"version"`
}

func (e *OrderEvent) String() string {
	return fmt.Sprintf("%s: [%s] %s %d", e.EventType, e.OrderID, e.Occured, e.Version)
}

// NewOrderEvent creates a new order event
func NewOrderEvent(eventType OrderEventType, orderID string, occured time.Time, version uint) *OrderEvent {
	return &OrderEvent{eventType, orderID, occured, version}
}
