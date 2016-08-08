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

// AddOrder adds a order to the price
func (op *OrderPrice) AddOrder(order *Order) {
	if order.Direction == Buy {
		op.Buy.Add(order)
	} else {
		op.Sell.Add(order)
	}
}
