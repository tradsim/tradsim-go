package models

import "github.com/satori/go.uuid"

// OrderBook contains all orders currently in the market
type OrderBook struct {
	Symbols map[string][]*OrderPrice
	Orders  map[uuid.UUID]*Order
}

// NewOrderBook creates a new order book
func NewOrderBook() *OrderBook {

	return &OrderBook{make(map[string][]*OrderPrice), make(map[uuid.UUID]*Order)}
}
