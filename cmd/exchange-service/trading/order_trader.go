package trading

import (
	"time"

	"github.com/mantzas/adaptlog"
	"github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

// Trader interface
type Trader interface {
	Trade(book *models.OrderBook, order *models.Order) (*models.Order, error)
}

// OrderTrader implementation
type OrderTrader struct {
	logger    adaptlog.LevelLogger
	publisher events.EventPublisher
}

// NewOrderTrader creates a new order trader
func NewOrderTrader(publisher events.EventPublisher) *OrderTrader {
	return &OrderTrader{adaptlog.NewStdLevelLogger("OrderTrader"), publisher}
}

// Trade processes a order against the book
func (ot *OrderTrader) Trade(book *models.OrderBook, order *models.Order) {

	prices, ok := book.Symbols[order.Symbol]
	if !ok {
		ot.logger.Debugf("Symbol %s not in book", order.Symbol)
		return
	}

	if order.Direction == models.Buy {
		ot.tradePricesBuy(prices, order)
	} else {
		ot.tradePricesSell(prices, order)
	}
}

func (ot *OrderTrader) tradePricesBuy(prices []models.OrderPrice, order *models.Order) {
	for _, price := range prices {
		if price.Price > order.Price {
			return
		}
		ot.tradePrice(&price, order)
	}
}

func (ot *OrderTrader) tradePricesSell(prices []models.OrderPrice, order *models.Order) {

	for i := len(prices) - 1; i >= 0; i-- {
		if prices[i].Price < order.Price {
			return
		}
		ot.tradePrice(&prices[i], order)
	}
}

func (ot *OrderTrader) tradePrice(price *models.OrderPrice, order *models.Order) {

	if order.Direction == models.Buy {
		if price.Buy.Quantity > 0 {

		}
	} else {
		if price.Sell.Quantity > 0 {

		}
	}
}

func (ot *OrderTrader) trade(existing *models.Order, new *models.Order) {
	traded := uint(0)

	if existing.Remaining() >= new.Quantity {
		traded = new.Quantity
	} else {
		traded = existing.Remaining()
	}

	existing.Trade(traded)
	ot.publishTradedEvent(existing.ID, existing.Price, traded)
	new.Trade(traded)
	ot.publishTradedEvent(new.ID, existing.Price, traded)
}

func (ot *OrderTrader) publishTradedEvent(ID uuid.UUID, price float64, traded uint) {

	ev := events.NewOrderTraded(ID.String(), time.Now().UTC(), price, traded, uint(1))
	env, err := events.NewOrderEventEnvelope(ev, ev.EventType)

	if err != nil {
		ot.logger.Errorf("Failed to create envelope: %s", err.Error())
	}

	err = ot.publisher.Publish(env)
	if err != nil {
		ot.logger.Errorf("Failed to publish traded event: %s", ev.String())
	}
}
