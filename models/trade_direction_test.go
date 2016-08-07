package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var tradeDirectionTests = []struct {
	in  TradeDirection
	out string
}{
	{Buy, "Buy"},
	{Sell, "Sell"},
	{9, "Not mapped value 9"},
}

func TestTradeDirectionString(t *testing.T) {

	for _, tt := range tradeDirectionTests {

		require.Equal(t, tt.in.String(), tt.out, "Expected %s but got %s", tt.out, tt.in.String())
	}
}

func TestTradeDirectionFromString(t *testing.T) {

	var tradeDirectionTests = []struct {
		in  string
		out TradeDirection
		err error
	}{
		{"Buy", Buy, nil},
		{"Sell", Sell, nil},
		{"9", 9, errors.New("Not mapped 9")},
	}

	require := require.New(t)

	for _, tt := range tradeDirectionTests {

		direction, err := TradeDirectionFromString(tt.in)

		if tt.err != nil {
			require.NotNil(err)
			require.Equal(err, tt.err)
		} else {
			require.Nil(err)
			require.Equal(direction, tt.out)
		}
	}
}
