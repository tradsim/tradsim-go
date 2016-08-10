package trading

import (
	"errors"

	"github.com/mantzas/adaptlog"
	"github.com/tradsim/tradsim-go/models"
)

// Appender interface
type Appender interface {
	Append(book *models.OrderBook, order *models.Order) error
}

// OrderAppender adds order to the book
type OrderAppender struct {
	logger adaptlog.LevelLogger
}

// NewOrderAppender creates a new order appender
func NewOrderAppender() *OrderAppender {
	return &OrderAppender{adaptlog.NewStdLevelLogger("OrderAppender")}
}

// Append the order to the book
func (oa *OrderAppender) Append(book *models.OrderBook, order *models.Order) error {

	if order.Status != models.Pending {

		return errors.New("Order status is not pending")
	}

	prices, ok := book.Symbols[order.Symbol]

	if !ok {
		oa.logger.Debugf("Symbol %s not found, adding symbol", order.Symbol)
		oa.addNewSymbol(book, order)
		return nil
	}
	oa.logger.Debugf("Symbol %s found", order.Symbol)

	found, i := findPrice(prices, order.Price)

	if found {
		oa.logger.Debugf("Price %f found, adding to price", order.Price)
		addOrderToPrice(prices[i], order)
	} else {
		oa.logger.Debugf("Price %f not found", order.Price)
		price := models.NewOrderPrice(order.Price)
		addOrderToPrice(price, order)

		if i == len(prices) { //append to the end
			prices = append(prices, price)
			oa.logger.Debugf("Price %f appended to the end", order.Price)
		} else { // append at i
			oa.logger.Debugf("Price %f appended at %d", order.Price, i)
			if i == 0 {
				prices = append([]*models.OrderPrice{price}, prices...)
			} else {
				prices = append(prices[:i], append([]*models.OrderPrice{price}, prices[i:]...)...)
			}
		}
		book.Symbols[order.Symbol] = prices
	}

	return nil
}

// this function returns the following
// found == true and index when price found
// found == false and at which index should it be inserted
func findPrice(prices []*models.OrderPrice, price float64) (bool, int) {

	for i := 0; i < len(prices); i++ {
		if prices[i].Price == price {
			return true, i
		} else if prices[i].Price > price {
			return false, i
		}
	}

	return false, len(prices)
}

func (oa *OrderAppender) addNewSymbol(book *models.OrderBook, order *models.Order) {

	price := models.NewOrderPrice(order.Price)
	addOrderToPrice(price, order)
	book.Symbols[order.Symbol] = []*models.OrderPrice{price}
	oa.logger.Debugf("addNewSymbol: Symbol %s appended", order.Symbol)
}

func addOrderToPrice(price *models.OrderPrice, order *models.Order) {
	if order.Direction == models.Buy {
		price.Buy.Quantity += order.Remaining()
		price.Buy.Orders = append(price.Buy.Orders, order)
	} else {
		price.Sell.Quantity += order.Remaining()
		price.Sell.Orders = append(price.Buy.Orders, order)
	}
}
