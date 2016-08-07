package models

// OrderQuantity defines the quantity and the orders underneath
type OrderQuantity struct {
	Quantity uint
	Orders   []Order
}

// NewOrderQuantity returns a new order quantity
func NewOrderQuantity() *OrderQuantity {
	return &OrderQuantity{0, make([]Order, 0)}
}
