package trading

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/models"
)

func TestAppendInvalidORder(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	ap := NewOrderAppender()

	order := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)
	order.Status = models.FullyFilled

	err := ap.Append(book, order)

	require.NotNil(err)
}

func TestAppendNewSymbol(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	ap := NewOrderAppender()

	order := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)

	ap.Append(book, order)

	require.Len(book.Symbols, 1)
	_, ok := book.Symbols["TT"]
	require.True(ok)
}

func TestAppendFoundOrder(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	ap := NewOrderAppender()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)
	order2 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 12, models.Sell)

	err := ap.Append(book, order1)
	require.Nil(err)
	err = ap.Append(book, order2)
	require.Nil(err)

	require.Len(book.Symbols, 1)
	prices, ok := book.Symbols["TT"]
	require.True(ok)
	require.Len(prices, 1)
	require.Equal(uint(22), prices[0].Sell.Quantity)
	require.Len(prices[0].Sell.Orders, 1)
}

func TestAppendAtTheEnd(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	ap := NewOrderAppender()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 199.98, 10, models.Sell)

	err := ap.Append(book, order1)

	require.Nil(err)
	require.Len(book.Symbols, 1)
	prices, ok := book.Symbols["TT"]
	require.True(ok)
	require.Len(prices, 1)
	require.Equal(199.98, prices[0].Price)
	require.Len(prices[0].Sell.Orders, 1)
	require.Equal(uint(10), prices[0].Sell.Quantity)

	order2 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 12, models.Sell)

	err = ap.Append(book, order2)

	require.Nil(err)
	require.Len(book.Symbols, 1)
	prices, ok = book.Symbols["TT"]
	require.True(ok)
	require.Len(prices, 2)
	require.Equal(199.99, prices[1].Price)
	require.Len(prices[1].Sell.Orders, 1)
	require.Equal(uint(12), prices[1].Sell.Quantity)
}

func TestAppendAtTheBeginning(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	ap := NewOrderAppender()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)
	order2 := models.NewOrder(uuid.NewV4(), "TT", 199.98, 12, models.Sell)

	err := ap.Append(book, order1)
	require.Nil(err)

	err = ap.Append(book, order2)
	require.Nil(err)

	require.Len(book.Symbols, 1)
	prices, ok := book.Symbols["TT"]
	require.True(ok)
	require.Len(prices, 2)
	require.Equal(199.98, prices[0].Price)
	require.Equal(199.99, prices[1].Price)
}

func TestAppendAtTheMiddle(t *testing.T) {

	require := require.New(t)

	book := models.NewOrderBook()

	ap := NewOrderAppender()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 199.99, 10, models.Sell)
	order2 := models.NewOrder(uuid.NewV4(), "TT", 199.97, 12, models.Sell)
	order3 := models.NewOrder(uuid.NewV4(), "TT", 199.98, 21, models.Sell)

	err := ap.Append(book, order1)
	require.Nil(err)

	err = ap.Append(book, order2)
	require.Nil(err)

	err = ap.Append(book, order3)
	require.Nil(err)

	require.Len(book.Symbols, 1)
	prices, ok := book.Symbols["TT"]
	require.True(ok)
	require.Len(prices, 3)
	require.Equal(199.97, prices[0].Price, "1")
	require.Equal(199.98, prices[1].Price, "2")
	require.Equal(199.99, prices[2].Price, "3")
}

// append
