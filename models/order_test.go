package models

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewOrder(t *testing.T) {

	require := require.New(t)

	u, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")
	order := NewOrder(u, "TT", 10.0, 10, Sell)

	require.NotNil(order)
}

func TestOrderString(t *testing.T) {

	require := require.New(t)

	order := getOrder(10, 0)

	require.Equal("[d1de4242-6620-4030-b2a7-4a701631c3ba] TT@199.990000 Buy 10/0/10 Pending", order.String())
}

func TestRemaining(t *testing.T) {

	require := require.New(t)

	order := getOrder(10, 2)

	require.Equal(uint(8), order.Remaining())
}

func TestAmend(t *testing.T) {

	require := require.New(t)

	order := getOrder(10, 0)
	order.Amend(2)

	require.Equal(uint(12), order.Quantity)
	require.Equal(Pending, order.Status)
}

func TestTrade(t *testing.T) {

	require := require.New(t)

	order := getOrder(10, 0)
	order.Trade(2)

	require.Equal(uint(10), order.Quantity)
	require.Equal(uint(2), order.Traded)
	require.Equal(PartiallyFilled, order.Status)
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

	require := require.New(t)

	for _, c := range cases {
		o := getOrder(2, c.in)
		o.UpdateStatus()
		require.Equal(c.out, o.Status)
	}
}

func getOrder(quantity uint, traded uint) *Order {

	u, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")

	return NewOrderFull(u, "TT", 199.99, quantity, traded, Buy, Pending)
}
