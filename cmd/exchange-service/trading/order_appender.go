package trading

import (
	"errors"
	"log"
	"sync"

	"github.com/tradsim/tradsim-go/models"
)

// Appender interface
type Appender interface {
	Append(book *models.OrderBook, order *models.Order) error
}

// OrderAppender adds order to the book
type OrderAppender struct {
	mu sync.Mutex
}

// NewOrderAppender creates a new order appender
func NewOrderAppender() *OrderAppender {
	return &OrderAppender{sync.Mutex{}}
}

// Append the order to the book
func (oa *OrderAppender) Append(book *models.OrderBook, order *models.Order) error {

	if order.Status != models.Pending {
		return errors.New("Order status is not pending")
	}

	oa.mu.Lock()
	defer oa.mu.Unlock()

	prices, ok := book.Symbols[order.Symbol]

	if !ok {
		oa.addNewSymbol(book, order)
		book.Orders[order.ID] = order
		return nil
	}

	found, i := findPrice(prices, order.Price)

	if found {
		log.Printf("Price %f found, adding to price", order.Price)
		addOrderToPrice(prices[i], order)
	} else {
		log.Printf("Price %f not found", order.Price)
		price := models.NewOrderPrice(order.Price)
		addOrderToPrice(price, order)

		if i == len(prices) { //append to the end
			prices = append(prices, price)
			log.Printf("Price %f appended to the end", order.Price)
		} else { // append at i
			log.Printf("Price %f appended at %d", order.Price, i)
			if i == 0 {
				prices = append([]*models.OrderPrice{price}, prices...)
			} else {
				prices = append(prices[:i], append([]*models.OrderPrice{price}, prices[i:]...)...)
			}
		}
		book.Symbols[order.Symbol] = prices
	}
	book.Orders[order.ID] = order

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
	log.Printf("Symbol %s appended", order.Symbol)
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
