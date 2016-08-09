package models

// OrderBook contains all orders currently in the market
type OrderBook struct {
	Symbols map[string][]*OrderPrice
}

// NewOrderBook creates a new order book
func NewOrderBook() *OrderBook {

	return &OrderBook{make(map[string][]*OrderPrice)}
}
