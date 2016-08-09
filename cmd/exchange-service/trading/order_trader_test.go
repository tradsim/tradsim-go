package trading

import (
	"os"
	"testing"

	"github.com/mantzas/adaptlog"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

func TestMain(m *testing.M) {
	adaptlog.ConfigureStdLevelLogger(adaptlog.DebugLevel, nil, "")
	retCode := m.Run()
	os.Exit(retCode)
}

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

func TestTradeBuyLeavingFullyFilled(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()
	orderSell1 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)
	orderSell2 := models.NewOrder(uuid.NewV4(), "TT", 199.98, 10, models.Sell)
	orderSell3 := models.NewOrder(uuid.NewV4(), "TT", 199.97, 10, models.Sell)

	ap := NewOrderAppender()
	ap.Append(book, orderSell1)
	ap.Append(book, orderSell2)
	ap.Append(book, orderSell3)

	orderBuy := models.NewOrder(uuid.NewV4(), "TT", 199.98, 20, models.Buy)

	publisher := &MockPublisher{}
	trader := NewOrderTrader(publisher)
	trader.Trade(book, orderBuy)

	prices := book.Symbols["TT"]

	require.Len(prices, 3)
	require.Equal(uint(20), orderBuy.Traded)
	require.Equal(models.FullyFilled, orderBuy.Status)
	require.Equal(uint(0), orderBuy.Remaining())
	require.Equal(uint(0), prices[0].Sell.Quantity, "Sell %f quantity %d", prices[0].Price, prices[0].Sell.Quantity)
	require.Equal(uint(0), prices[1].Sell.Quantity, "Sell %f quantity %d", prices[1].Price, prices[1].Sell.Quantity)
	require.Equal(uint(10), prices[2].Sell.Quantity, "Sell %f quantity %d", prices[2].Price, prices[2].Sell.Quantity)

	require.Len(publisher.envelopes, 4)
}

func TestTradeSellLeavingFullyFilled(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()
	orderBuy1 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Buy)
	orderBuy2 := models.NewOrder(uuid.NewV4(), "TT", 199.98, 10, models.Buy)
	orderBuy3 := models.NewOrder(uuid.NewV4(), "TT", 199.97, 10, models.Buy)

	ap := NewOrderAppender()
	ap.Append(book, orderBuy1)
	ap.Append(book, orderBuy2)
	ap.Append(book, orderBuy3)

	orderSell := models.NewOrder(uuid.NewV4(), "TT", 199.98, 20, models.Sell)

	publisher := &MockPublisher{}
	trader := NewOrderTrader(publisher)
	trader.Trade(book, orderSell)

	prices := book.Symbols["TT"]

	require.Len(prices, 3)
	require.Equal(uint(20), orderSell.Traded)
	require.Equal(models.FullyFilled, orderSell.Status)
	require.Equal(uint(0), orderSell.Remaining())
	require.Equal(uint(10), prices[0].Buy.Quantity, "Sell %f quantity %d", prices[0].Price, prices[0].Buy.Quantity)
	require.Equal(uint(0), prices[1].Buy.Quantity, "Sell %f quantity %d", prices[1].Price, prices[1].Buy.Quantity)
	require.Equal(uint(0), prices[2].Buy.Quantity, "Sell %f quantity %d", prices[2].Price, prices[2].Buy.Quantity)
	require.Len(publisher.envelopes, 4)
}

type MockPublisher struct {
	envelopes []*events.OrderEventEnvelope
}

func (mp *MockPublisher) Open() error {
	return nil
}

func (mp *MockPublisher) Close() {
}

func (mp *MockPublisher) Publish(envelope *events.OrderEventEnvelope) error {
	mp.envelopes = append(mp.envelopes, envelope)
	return nil
}
