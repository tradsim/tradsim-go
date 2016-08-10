package models

import (
	"testing"

	"os"

	"github.com/mantzas/adaptlog"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	adaptlog.ConfigureStdLevelLogger(adaptlog.DebugLevel, nil, "")
	retCode := m.Run()
	os.Exit(retCode)
}

func TestNewOrderPrice(t *testing.T) {

	require := require.New(t)

	price := NewOrderPrice(199.99)

	require.Equal(199.99, price.Price)
	require.Equal(uint(0), price.Sell.Quantity)
	require.Len(price.Sell.Orders, 0)
	require.Equal(uint(0), price.Buy.Quantity)
	require.Len(price.Buy.Orders, 0)
}
