package trading

import (
	"log"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

// Trader interface
type Trader interface {
	Trade(book *models.OrderBook, order *models.Order)
}

// OrderTrader implementation
type OrderTrader struct {
	publisher events.EventPublisher
	mu        sync.Mutex
}

// NewOrderTrader creates a new order trader
func NewOrderTrader(publisher events.EventPublisher) *OrderTrader {
	return &OrderTrader{publisher, sync.Mutex{}}
}

// Trade processes a order against the book
func (ot *OrderTrader) Trade(book *models.OrderBook, order *models.Order) {

	ot.mu.Lock()
	defer ot.mu.Unlock()

	prices, ok := book.Symbols[order.Symbol]
	if !ok {
		log.Printf("Symbol %s not in book", order.Symbol)
		return
	}

	if order.Direction == models.Buy {
		ot.tradePricesBuy(prices, order)
	} else {
		ot.tradePricesSell(prices, order)
	}
}

func (ot *OrderTrader) tradePricesBuy(prices []*models.OrderPrice, order *models.Order) {
	for _, price := range prices {
		if price.Price > order.Price {
			log.Printf("Price %f is greater than order price %f", price.Price, order.Price)
			return
		}
		log.Printf("Trading with price %f. order price %f", price.Price, order.Price)
		ot.tradePrice(price, order)
	}
}

func (ot *OrderTrader) tradePricesSell(prices []*models.OrderPrice, order *models.Order) {

	for i := len(prices) - 1; i >= 0; i-- {
		if prices[i].Price < order.Price {
			log.Printf("Price %f is less than order price %f", prices[i].Price, order.Price)
			return
		}
		log.Printf("Trading with price %f. order price %f", prices[i].Price, order.Price)
		ot.tradePrice(prices[i], order)
	}
}

func (ot *OrderTrader) tradePrice(price *models.OrderPrice, order *models.Order) {

	if order.Direction == models.Buy {
		if price.Sell.Quantity > 0 {
			log.Printf("Sell quantity on price %f is %d and order count %d", price.Price, price.Sell.Quantity, len(price.Sell.Orders))
			for _, existing := range price.Sell.Orders {
				ot.trade(existing, order)
			}
			price.Sell.Quantity = ot.compactOrdersAndGetQuantity(&price.Sell.Orders)
			log.Printf("Sell quantity after trade on price %f is %d and order count %d", price.Price, price.Sell.Quantity, len(price.Sell.Orders))
		} else {
			log.Printf("Sell quantity on price %f is zero", price.Price)
		}
	} else {
		if price.Buy.Quantity > 0 {
			log.Printf("Buy quantity on price %f is %d and order count %d", price.Price, price.Buy.Quantity, len(price.Buy.Orders))
			for _, existing := range price.Buy.Orders {
				ot.trade(existing, order)
			}
			price.Buy.Quantity = ot.compactOrdersAndGetQuantity(&price.Buy.Orders)
			log.Printf("Buy quantity after trade on price %f is %d and order count %d", price.Price, price.Buy.Quantity, len(price.Buy.Orders))
		} else {
			log.Printf("Buy quantity on price %f is zero", price.Price)
		}
	}
}

func (ot *OrderTrader) trade(existing *models.Order, new *models.Order) {

	if !existing.Status.IsTradeable() {
		return
	}

	traded := uint(0)

	if existing.Remaining() >= new.Remaining() {
		traded = new.Remaining()
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
		log.Printf("Failed to create envelope: %s", err.Error())
	}

	err = ot.publisher.Publish(env)
	if err != nil {
		log.Printf("Failed to publish traded event: %s", ev.String())
	}
}

func (ot *OrderTrader) compactOrdersAndGetQuantity(orders *[]*models.Order) uint {
	quantity := uint(0)
	var temp = make([]*models.Order, 0)
	for _, order := range *orders {
		if order.Status != models.FullyFilled && order.Status != models.OverFilled &&
			order.Status != models.Cancelled {
			quantity += order.Remaining()
			temp = append(temp, order)
		}
	}
	*orders = temp
	return quantity
}
