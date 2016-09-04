package aggregator

import (
	"testing"
	"time"

	incmodel "github.com/mantzas/incata/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

func TestAggregationSuccess(t *testing.T) {

	require := require.New(t)
	orderID, err := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")
	accepted := events.NewOrderAccepted(orderID.String(), time.Now().UTC(), "TT", 1.99, 10, models.Buy, 1)
	traded := events.NewOrderTraded(orderID.String(), time.Now().UTC(), 1.98, 10, 1)
	amended := events.NewOrderAmended(orderID.String(), 20, time.Now().UTC(), 1)
	cancelled := events.NewOrderCancelled(orderID.String(), time.Now().UTC(), 1)

	evs := []incmodel.Event{*incmodel.NewEvent(orderID, time.Now().UTC(), *accepted, string(accepted.EventType), 1),
		*incmodel.NewEvent(orderID, time.Now().UTC(), *traded, string(traded.EventType), 1),
		*incmodel.NewEvent(orderID, time.Now().UTC(), *amended, string(amended.EventType), 1),
		*incmodel.NewEvent(orderID, time.Now().UTC(), *cancelled, string(cancelled.EventType), 1)}

	ag := NewEventAggregator()

	o, err := ag.Aggregate(evs)

	require.Nil(err)
	require.Equal(models.Cancelled, o.Status)
	require.Equal(uint(20), o.Quantity)
	require.Equal(uint(10), o.Traded)
	require.Len(o.Logs, 4)
	require.Equal(string(events.OrderAcceptedType), o.Logs[0].Action)
	require.Equal(string(events.OrderTradedType), o.Logs[1].Action)
	require.Equal(string(events.OrderAmendedType), o.Logs[2].Action)
	require.Equal(string(events.OrderCancelledType), o.Logs[3].Action)
	require.Len(o.Trades, 1)
}

func TestAggregationInvalidEventSuccess(t *testing.T) {

	require := require.New(t)
	orderID, err := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")
	stored := events.NewOrderEventStored(orderID.String(), time.Now().UTC(), 1)

	evs := []incmodel.Event{*incmodel.NewEvent(orderID, time.Now().UTC(), *stored, string(stored.EventType), 1)}

	ag := NewEventAggregator()

	_, err = ag.Aggregate(evs)

	require.NotNil(err)
}
