package models

import (
	"testing"

	"github.com/mantzas/adaptlog"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"os"
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

func TestAddOrder(t *testing.T) {

	require := require.New(t)

	price := NewOrderPrice(199.99)
	price.AddOrder(NewOrder(uuid.NewV4(), "TT", 199.99, 10, Sell))
	price.AddOrder(NewOrder(uuid.NewV4(), "TT", 199.99, 1, Buy))
	price.AddOrder(NewOrder(uuid.NewV4(), "TT", 199.99, 2, Sell))

	require.Equal(199.99, price.Price)
	require.Equal(uint(12), price.Sell.Quantity)
	require.Len(price.Sell.Orders, 2)
	require.Equal(uint(1), price.Buy.Quantity)
	require.Len(price.Buy.Orders, 1)
}
