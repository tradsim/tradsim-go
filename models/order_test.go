package models

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrderString(t *testing.T) {

	assert := assert.New(t)

	order := getOrder(10, 10)

	assert.Equal("[d1de4242-6620-4030-b2a7-4a701631c3ba] TT@199.990000 Buy 10 Pending", order.String())
}

func TestOrderSetStatusPending(t *testing.T) {

	assert := assert.New(t)

	order := getOrder(10, 0)
	order.SetStatus()

	assert.Equal(order.Status, Pending)
}

func TestOrderSetStatusPartiallyFilled(t *testing.T) {

	assert := assert.New(t)

	order := getOrder(10, 5)
	order.SetStatus()

	assert.Equal(order.Status, PartiallyFilled)
}

func TestOrderSetStatusFullyFilled(t *testing.T) {

	assert := assert.New(t)

	order := getOrder(10, 10)
	order.SetStatus()

	assert.Equal(order.Status, FullyFilled)
}

func TestOrderSetStatusOverFilled(t *testing.T) {

	assert := assert.New(t)

	order := getOrder(10, 12)
	order.SetStatus()

	assert.Equal(order.Status, OverFilled)
}

func getOrder(quantity uint, traded uint) *Order {

	u, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")

	return NewOrderFull(u, "TT", 199.99, quantity, traded, Buy, Pending)
}

func getOrderFull(price float64, quantity uint, direction TradeDirection) *Order {

	u, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")

	return NewOrderFull(u, "TT", price, quantity, 0, direction, Pending)
}
