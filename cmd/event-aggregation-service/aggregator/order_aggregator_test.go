package aggregator

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/cmd/event-aggregation-service/models"
	commonmodel "github.com/tradsim/tradsim-go/models"
)

func TestOrderAggregator(t *testing.T) {

	require := require.New(t)
	orderID, _ := uuid.FromString("d1de4242-6620-4030-b2a7-4a701631c3ba")

	updated := time.Now().UTC()

	order1 := models.NewOrder(orderID, "TT", 1.99, 10, commonmodel.Buy, commonmodel.Pending, updated.Add(-1*time.Hour))
	order1.Trade(1.95, uint(10), updated.Add(-1*time.Hour))
	order2 := models.NewOrder(orderID, "TT", 1.99, 10, commonmodel.Sell, commonmodel.Pending, updated)
	order2.Trade(2.0, uint(5), updated)
	order3 := models.NewOrder(orderID, "ETE", 1.99, 5, commonmodel.Sell, commonmodel.Pending, updated)

	orders := []models.Order{*order1, *order2, *order3}

	ag := NewOrderAggregator()

	pos := ag.Aggregate("TT", orders)

	require.Equal("TT", pos.Symbol)
	require.Equal(updated, pos.Updated)
	require.Equal(5, pos.Quantity)
}
