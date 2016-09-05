package models

import (
	"time"

	"github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

// Order defines the order state
type Order struct {
	ID          uuid.UUID
	Symbol      string
	Price       float64
	Quantity    uint
	Traded      uint
	TradedPrice float64
	Direction   models.TradeDirection
	Status      models.OrderStatus
	Created     time.Time
	Updated     time.Time
	Trades      []Trade
	Logs        []OrderLog
}

// NewOrder creates a new order
func NewOrder(id uuid.UUID, symbol string, price float64, quantity uint, direction models.TradeDirection, status models.OrderStatus, created time.Time) *Order {
	o := Order{id, symbol, price, quantity, uint(0), 0.0, direction, status, created, created, make([]Trade, 0), make([]OrderLog, 0)}
	o.appendLog(string(events.OrderAcceptedType), created)
	return &o
}

// Trade appends trade to order and updates order
func (o *Order) Trade(p float64, q uint, t time.Time) {
	o.Traded += q
	o.Trades = append(o.Trades, Trade{0, o.ID, p, q, t})
	o.updateTradedPrice()
	o.Status = models.ResolveStatus(o.Quantity, o.Traded)
	o.Updated = t
	o.appendLog(string(events.OrderTradedType), t)
}

// Amend amends the order
func (o *Order) Amend(q uint, t time.Time) {
	o.Quantity = q
	o.Status = models.ResolveStatus(o.Quantity, o.Traded)
	o.Updated = t
	o.appendLog(string(events.OrderAmendedType), t)
}

// Cancel the order
func (o *Order) Cancel(t time.Time) {
	o.Status = models.Cancelled
	o.Updated = t
	o.appendLog(string(events.OrderCancelledType), t)
}

func (o *Order) appendLog(a string, t time.Time) {
	o.Logs = append(o.Logs, OrderLog{0, o.ID, a, t})
}

func (o *Order) updateTradedPrice() {

	price := 0.0
	quantity := uint(0)

	for _, trade := range o.Trades {
		price += float64(trade.Quantity) * trade.Price
		quantity += trade.Quantity
	}

	if quantity == uint(0) {
		o.TradedPrice = 0.0
	} else {
		o.TradedPrice = price / float64(quantity)
	}
}
