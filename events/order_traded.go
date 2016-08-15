package events

import (
	"fmt"
	"time"
)

// OrderTraded defines a order traded event
type OrderTraded struct {
	OrderEvent
	Price    float64 `json:"price"`
	Quantity uint    `json:"quantity"`
}

func (e *OrderTraded) String() string {
	return fmt.Sprintf("%s %d@%f", e.OrderEvent.String(), e.Quantity, e.Price)
}

// NewOrderTraded creates a new order traded event
func NewOrderTraded(orderID string, occured time.Time, price float64, quantity uint, version uint) *OrderTraded {

	return &OrderTraded{*NewOrderEvent(OrderTradedType, orderID, occured, version), price, quantity}
}
