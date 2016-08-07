package models

// OrderPrice define the price and the buy and sell sides
type OrderPrice struct {
	Price float64
	Buy   OrderQuantity
	Sell  OrderQuantity
}

// NewOrderPrice creates a new order price
func NewOrderPrice(price float64) *OrderPrice {
	return &OrderPrice{price, *NewOrderQuantity(), *NewOrderQuantity()}
}
