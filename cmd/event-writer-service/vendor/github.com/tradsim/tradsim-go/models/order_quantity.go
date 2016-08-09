package models

// OrderQuantity defines the quantity and the orders underneath
type OrderQuantity struct {
	Quantity uint
	Orders   []*Order
}

// NewOrderQuantity returns a new order quantity
func NewOrderQuantity() *OrderQuantity {
	return &OrderQuantity{0, make([]*Order, 0)}
}

// Add order
func (oq *OrderQuantity) Add(order *Order) {

	oq.Quantity += order.Quantity
	oq.Orders = append(oq.Orders, order)
}

// Update quantity
func (oq *OrderQuantity) Update() {
	quantity := uint(0)

	for _, order := range oq.Orders {
		if order.Status != FullyFilled || order.Status != OverFilled {
			quantity += order.Remaining()
		}
	}

	oq.Quantity = quantity
}
