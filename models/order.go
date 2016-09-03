package models

import (
	"fmt"

	"github.com/satori/go.uuid"
)

// Order defines a order for a specific symbol
type Order struct {
	ID        uuid.UUID
	Symbol    string
	Price     float64
	Quantity  uint
	Traded    uint
	Direction TradeDirection
	Status    OrderStatus
}

// NewOrder creates a new order
func NewOrder(id uuid.UUID, symbol string, price float64, quantity uint, direction TradeDirection) *Order {
	return &Order{id, symbol, price, quantity, 0, direction, Pending}
}

// NewOrderFull creates a new order with all parameters
func NewOrderFull(id uuid.UUID, symbol string, price float64, quantity uint, traded uint, direction TradeDirection, orderStatus OrderStatus) *Order {
	return &Order{id, symbol, price, quantity, traded, direction, orderStatus}
}

// Remaining return the reamining quantity
func (o *Order) Remaining() uint {
	return o.Quantity - o.Traded
}

// Amend order by quantity
func (o *Order) Amend(q uint) {
	o.Quantity += q
	o.UpdateStatus()
}

// Trade order
func (o *Order) Trade(q uint) {
	o.Traded += q
	o.UpdateStatus()
}

func (o *Order) String() string {
	return fmt.Sprintf("[%s] %s@%f %s %d/%d/%d %s", o.ID, o.Symbol, o.Price, o.Direction.String(), o.Quantity, o.Traded, o.Remaining(), o.Status.String())
}

// UpdateStatus updates the status based on the quantities
func (o *Order) UpdateStatus() {
	if o.Traded == uint(0) {
		o.Status = Pending
	} else if o.Traded < o.Quantity {
		o.Status = PartiallyFilled
	} else if o.Traded == o.Quantity {
		o.Status = FullyFilled
	} else {
		o.Status = OverFilled
	}
}
