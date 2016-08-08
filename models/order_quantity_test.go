package models

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewOrderQuantity(t *testing.T) {

	require := require.New(t)

	quantity := NewOrderQuantity()

	require.Equal(uint(0), quantity.Quantity)
	require.Len(quantity.Orders, 0)
}

func TestAdd(t *testing.T) {

	require := require.New(t)

	quantity := NewOrderQuantity()
	quantity.Add(NewOrder(uuid.NewV4(), "TT", 199.99, 10, Sell))

	require.Equal(uint(10), quantity.Quantity)
	require.Len(quantity.Orders, 1)
}
