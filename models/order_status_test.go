package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var orderStatusTests = []struct {
	in  OrderStatus
	out string
}{
	{Pending, "Pending"},
	{PartiallyFilled, "PartiallyFilled"},
	{FullyFilled, "FullyFilled"},
	{OverFilled, "OverFilled"},
	{9, "Not mapped value 9"},
}

func TestOrderStatusString(t *testing.T) {

	for _, tt := range orderStatusTests {

		require.Equal(t, tt.in.String(), tt.out, "Expected %s but got %s", tt.out, tt.in.String())
	}
}
