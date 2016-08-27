package trading

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/mocks"
	"github.com/tradsim/tradsim-go/models"
)

func TestAmendMissingSymbolReturnsFalse(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	order := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)

	am := NewOrderAmender(&mocks.MockPublisher{})

	require.False(am.Amend(book, order))
}

func TestAmendMissingPriceReturnsFalse(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)

	ap := NewOrderAppender()
	ap.Append(book, order1)

	am := NewOrderAmender(&mocks.MockPublisher{})

	order2 := models.NewOrder(uuid.NewV4(), "TT", 299.99, 10, models.Sell)

	require.False(am.Amend(book, order2))
}

func TestAmendMissingOrderReturnsFalse(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)

	ap := NewOrderAppender()
	ap.Append(book, order1)

	am := NewOrderAmender(&mocks.MockPublisher{})

	order2 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)

	require.False(am.Amend(book, order2))
}

func TestAmendBuy(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Buy)

	ap := NewOrderAppender()
	ap.Append(book, order1)

	am := NewOrderAmender(&mocks.MockPublisher{})

	order2 := models.NewOrder(order1.ID, "TT", 199.99, 20, models.Buy)

	prices, _ := book.Symbols["TT"]

	require.True(am.Amend(book, order2))
	require.Len(prices, 1)
	require.Equal(uint(20), prices[0].Buy.Quantity)
	require.Equal(uint(20), prices[0].Buy.Orders[0].Quantity)
}

func TestAmendSell(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)

	ap := NewOrderAppender()
	ap.Append(book, order1)

	am := NewOrderAmender(&mocks.MockPublisher{})

	order2 := models.NewOrder(order1.ID, "TT", 199.99, 20, models.Sell)

	prices, _ := book.Symbols["TT"]

	require.True(am.Amend(book, order2))
	require.Len(prices, 1)
	require.Equal(uint(20), prices[0].Sell.Quantity)
	require.Equal(uint(20), prices[0].Sell.Orders[0].Quantity)
}
