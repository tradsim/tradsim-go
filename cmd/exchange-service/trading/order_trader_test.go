package trading

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

func TestNewOrderTrader(t *testing.T) {

	require := require.New(t)

	trader := NewOrderTrader(&MockPublisher{})

	require.NotNil(trader)
}

func TestTradeNoSymbol(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()
	order := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)

	trader := NewOrderTrader(&MockPublisher{})
	trader.Trade(book, order)

	require.Equal(uint(0), order.Traded)
}

type MockPublisher struct {
}

func (mp *MockPublisher) Open() error {
	return nil
}

func (mp *MockPublisher) Close() {
}

func (mp *MockPublisher) Publish(envelop *events.OrderEventEnvelope) error {
	return nil
}
