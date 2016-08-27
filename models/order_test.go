package models

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {

	assert := assert.New(t)

	u, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")
	order := NewOrder(u, "TT", 10.0, 10, Sell)

	assert.NotNil(order)
}

func TestOrderString(t *testing.T) {

	assert := assert.New(t)

	order := getOrder(10, 0)

	assert.Equal("[d1de4242-6620-4030-b2a7-4a701631c3ba] TT@199.990000 Buy 10/0/10 Pending", order.String())
}

func TestRemaining(t *testing.T) {

	assert := assert.New(t)

	order := getOrder(10, 2)

	assert.Equal(uint(8), order.Remaining())
}

func TestAmend(t *testing.T) {

	assert := assert.New(t)

	order := getOrder(10, 0)
	order.Amend(2)

	assert.Equal(uint(12), order.Quantity)
	assert.Equal(Pending, order.Status)
}

func TestTrade(t *testing.T) {

	assert := assert.New(t)

	order := getOrder(10, 0)
	order.Trade(2)

	assert.Equal(uint(10), order.Quantity)
	assert.Equal(uint(2), order.Traded)
	assert.Equal(PartiallyFilled, order.Status)
}

func TestSetStatus(t *testing.T) {

	var cases = []struct {
		in  uint
		out OrderStatus
	}{
		{0, Pending},
		{1, PartiallyFilled},
		{2, FullyFilled},
		{3, OverFilled},
	}

	assert := assert.New(t)

	for _, c := range cases {
		o := getOrder(2, c.in)
		o.setStatus()
		assert.Equal(c.out, o.Status)
	}
}

func getOrder(quantity uint, traded uint) *Order {

	u, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")

	return NewOrderFull(u, "TT", 199.99, quantity, traded, Buy, Pending)
}
