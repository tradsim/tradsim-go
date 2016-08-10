package models

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

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

func getOrder(quantity uint, traded uint) *Order {

	u, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")

	return NewOrderFull(u, "TT", 199.99, quantity, traded, Buy, Pending)
}
