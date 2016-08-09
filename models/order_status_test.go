package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var orderStatusTradeableTests = []struct {
	in  OrderStatus
	out bool
}{
	{Pending, true},
	{PartiallyFilled, true},
	{FullyFilled, false},
	{OverFilled, false},
}

func TestOrderStatusIsTradeable(t *testing.T) {

	for _, tt := range orderStatusTradeableTests {

		require.Equal(t, tt.out, tt.in.IsTradeable(), "Expected %b but got %b", tt.out, tt.in.IsTradeable())
	}
}

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

var orderStatusFromStringTests = []struct {
	in  string
	out struct {
		status   OrderStatus
		hasError bool
	}
}{
	{"Pending", struct {
		status   OrderStatus
		hasError bool
	}{Pending, false}},
	{"PartiallyFilled", struct {
		status   OrderStatus
		hasError bool
	}{PartiallyFilled, false}},
	{"FullyFilled", struct {
		status   OrderStatus
		hasError bool
	}{FullyFilled, false}},
	{"OverFilled", struct {
		status   OrderStatus
		hasError bool
	}{OverFilled, false}},
	{"9", struct {
		status   OrderStatus
		hasError bool
	}{OverFilled, true}},
}

func TestOrderStatusFromString(t *testing.T) {

	for _, tt := range orderStatusFromStringTests {

		status, err := OrderStatusFromString(tt.in)

		require := require.New(t)

		if tt.out.hasError {
			require.NotNil(err)
		} else {
			require.Nil(err)
			require.Equal(tt.out.status, status, "Expected %v but got %v", tt.out.status, status)
		}
	}
}
