package models

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

func TestNewOrder(t *testing.T) {

	require := require.New(t)
	orderID, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")
	created := time.Now().UTC()

	o := NewOrder(orderID, "TT", 1.99, uint(10), models.Buy, models.Pending, created)

	require.Equal(orderID, o.ID)
	require.Equal("TT", o.Symbol)
	require.Equal(1.99, o.Price)
	require.Equal(uint(10), o.Quantity)
	require.Equal(uint(0), o.Traded)
	require.Equal(0.0, o.TradedPrice)
	require.Equal(models.Buy, o.Direction)
	require.Equal(models.Pending, o.Status)
	require.Equal(created, o.Created)
	require.Equal(created, o.Updated)
	require.Len(o.Trades, 0)
	require.Len(o.Logs, 1)
	require.Equal(string(events.OrderAcceptedType), o.Logs[0].Action)
	require.Equal(created, o.Logs[0].Occured)
}

func TestAmend(t *testing.T) {

	require := require.New(t)
	orderID, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")
	created := time.Now().UTC()
	updated := time.Now().UTC()
	o := NewOrder(orderID, "TT", 1.99, uint(10), models.Buy, models.Pending, created)
	o.Amend(20, updated)

	require.Equal(orderID, o.ID)
	require.Equal("TT", o.Symbol)
	require.Equal(1.99, o.Price)
	require.Equal(uint(20), o.Quantity)
	require.Equal(uint(0), o.Traded)
	require.Equal(0.0, o.TradedPrice)
	require.Equal(models.Buy, o.Direction)
	require.Equal(models.Pending, o.Status)
	require.Equal(created, o.Created)
	require.Equal(updated, o.Updated)
	require.Len(o.Trades, 0)
	require.Len(o.Logs, 2)
	require.Equal(string(events.OrderAcceptedType), o.Logs[0].Action)
	require.Equal(created, o.Logs[0].Occured)
	require.Equal(string(events.OrderAmendedType), o.Logs[1].Action)
	require.Equal(updated, o.Logs[1].Occured)
}

func TestCancel(t *testing.T) {

	require := require.New(t)
	orderID, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")
	created := time.Now().UTC()
	updated := time.Now().UTC()
	o := NewOrder(orderID, "TT", 1.99, uint(10), models.Buy, models.Pending, created)
	o.Cancel(updated)

	require.Equal(orderID, o.ID)
	require.Equal("TT", o.Symbol)
	require.Equal(1.99, o.Price)
	require.Equal(uint(10), o.Quantity)
	require.Equal(uint(0), o.Traded)
	require.Equal(0.0, o.TradedPrice)
	require.Equal(models.Buy, o.Direction)
	require.Equal(models.Cancelled, o.Status)
	require.Equal(created, o.Created)
	require.Equal(updated, o.Updated)
	require.Len(o.Trades, 0)
	require.Len(o.Logs, 2)
	require.Equal(string(events.OrderAcceptedType), o.Logs[0].Action)
	require.Equal(created, o.Logs[0].Occured)
	require.Equal(string(events.OrderCancelledType), o.Logs[1].Action)
	require.Equal(updated, o.Logs[1].Occured)
}

func TestAppendTrade(t *testing.T) {

	require := require.New(t)
	orderID, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")
	created := time.Now().UTC()

	o := NewOrder(orderID, "TT", 1.99, uint(10), models.Buy, models.Pending, created)

	updated := time.Now().UTC()
	o.Trade(1.95, 5, time.Now().UTC())
	o.Trade(1.96, 2, updated)

	require.Equal(orderID, o.ID)
	require.Equal("TT", o.Symbol)
	require.Equal(1.99, o.Price)
	require.Equal(uint(10), o.Quantity)
	require.Equal(uint(7), o.Traded)
	require.Equal(1.9528571428571428, o.TradedPrice)
	require.Equal(models.Buy, o.Direction)
	require.Equal(models.PartiallyFilled, o.Status)
	require.Equal(created, o.Created)
	require.Equal(updated, o.Updated)
	require.Len(o.Trades, 2)
	require.Len(o.Logs, 3)
	require.Equal(string(events.OrderAcceptedType), o.Logs[0].Action)
	require.Equal(string(events.OrderTradedType), o.Logs[1].Action, "%v", o.Logs)
	require.Equal(string(events.OrderTradedType), o.Logs[2].Action)
}
