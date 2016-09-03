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
	{Cancelled, false},
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
	{Pending, PendingText},
	{PartiallyFilled, PartiallyFilledText},
	{FullyFilled, FullyFilledText},
	{OverFilled, OverFilledText},
	{Cancelled, CancelledText},
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
	{PendingText, struct {
		status   OrderStatus
		hasError bool
	}{Pending, false}},
	{PartiallyFilledText, struct {
		status   OrderStatus
		hasError bool
	}{PartiallyFilled, false}},
	{FullyFilledText, struct {
		status   OrderStatus
		hasError bool
	}{FullyFilled, false}},
	{OverFilledText, struct {
		status   OrderStatus
		hasError bool
	}{OverFilled, false}},
	{CancelledText, struct {
		status   OrderStatus
		hasError bool
	}{Cancelled, false}},
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

func TestResolveStatus(t *testing.T) {

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

		require.Equal(c.out, ResolveStatus(2, c.in))
	}
}
