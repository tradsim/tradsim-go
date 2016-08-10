package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOrderQuantity(t *testing.T) {

	require := require.New(t)

	quantity := NewOrderQuantity()

	require.Equal(uint(0), quantity.Quantity)
	require.Len(quantity.Orders, 0)
}
