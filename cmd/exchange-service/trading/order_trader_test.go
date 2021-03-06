package trading

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/mocks"
	"github.com/tradsim/tradsim-go/models"
)

func TestNewOrderTrader(t *testing.T) {

	require := require.New(t)

	trader := NewOrderTrader(&mocks.MockPublisher{})

	require.NotNil(trader)
}

func TestTradeNoSymbol(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()
	order := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)

	trader := NewOrderTrader(&mocks.MockPublisher{})
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

	publisher := &mocks.MockPublisher{}
	trader := NewOrderTrader(publisher)
	trader.Trade(book, orderBuy)

	prices := book.Symbols["TT"]

	require.Len(prices, 3)
	require.Equal(uint(20), orderBuy.Traded)
	require.Equal(models.FullyFilled, orderBuy.Status)
	require.Equal(uint(0), orderBuy.Remaining())
	require.Equal(uint(0), prices[0].Sell.Quantity, "Sell %f quantity %d", prices[0].Price, prices[0].Sell.Quantity)
	require.Len(prices[0].Sell.Orders, 0)
	require.Equal(uint(0), prices[1].Sell.Quantity, "Sell %f quantity %d", prices[1].Price, prices[1].Sell.Quantity)
	require.Len(prices[1].Sell.Orders, 0)
	require.Equal(uint(10), prices[2].Sell.Quantity, "Sell %f quantity %d", prices[2].Price, prices[2].Sell.Quantity)
	require.Len(prices[2].Sell.Orders, 1)

	require.Len(publisher.Envelopes, 4)
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

	publisher := &mocks.MockPublisher{}
	trader := NewOrderTrader(publisher)
	trader.Trade(book, orderSell)

	prices := book.Symbols["TT"]

	require.Len(prices, 3)
	require.Equal(uint(20), orderSell.Traded)
	require.Equal(models.FullyFilled, orderSell.Status)
	require.Equal(uint(0), orderSell.Remaining())
	require.Equal(uint(10), prices[0].Buy.Quantity, "Sell %f quantity %d", prices[0].Price, prices[0].Buy.Quantity)
	require.Len(prices[0].Buy.Orders, 1)
	require.Equal(uint(0), prices[1].Buy.Quantity, "Sell %f quantity %d", prices[1].Price, prices[1].Buy.Quantity)
	require.Len(prices[1].Buy.Orders, 0)
	require.Equal(uint(0), prices[2].Buy.Quantity, "Sell %f quantity %d", prices[2].Price, prices[2].Buy.Quantity)
	require.Len(prices[2].Buy.Orders, 0)
	require.Len(publisher.Envelopes, 4)
}
