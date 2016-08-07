package models

import (
	"fmt"

	"github.com/mantzas/adaptlog"
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
	logger    adaptlog.LevelLogger
}

// NewOrder creates a new order
func NewOrder(id uuid.UUID, symbol string, price float64, quantity uint, direction TradeDirection) *Order {
	order := Order{id, symbol, price, quantity, 0, direction, Pending, adaptlog.NewStdLevelLogger("Order")}
	order.SetStatus()
	return &order
}

// NewOrderFull creates a new order with all parameters
func NewOrderFull(id uuid.UUID, symbol string, price float64, quantity uint, traded uint, direction TradeDirection, orderStatus OrderStatus) *Order {
	return &Order{id, symbol, price, quantity, traded, direction, orderStatus, adaptlog.NewStdLevelLogger("Order")}
}

func (o *Order) String() string {
	return fmt.Sprintf("[%s] %s@%f %s %d %s", o.ID, o.Symbol, o.Price, o.Direction.String(), o.Quantity, o.Status.String())
}

// SetStatus based on the quantities
func (o *Order) SetStatus() {
	if o.Traded == 0 {
		o.Status = Pending
		o.logger.Debugf("SetStatus: [%s] Traded 0, set to Pending", o.ID)
	} else if o.Traded < o.Quantity {
		o.Status = PartiallyFilled
		o.logger.Debugf("SetStatus: [%s] Traded less than Quantity, set to Partially Filled", o.ID)
	} else if o.Traded == o.Quantity {
		o.Status = FullyFilled
		o.logger.Debugf("SetStatus: [%s] Traded equals Quantity, set to Fully Filled", o.ID)
	} else {
		o.Status = OverFilled
		o.logger.Debugf("SetStatus: [%s] Traded greater than Quantity, set to Over Filled", o.ID)
	}
}
