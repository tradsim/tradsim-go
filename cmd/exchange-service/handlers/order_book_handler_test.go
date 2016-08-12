package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"
	"io/ioutil"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/cmd/exchange-service/trading"
	"github.com/tradsim/tradsim-go/models"
)

func TestGetSymbolHandler(t *testing.T) {

	require := require.New(t)
	book := models.NewOrderBook()

	buyOrder := models.NewOrder(uuid.NewV4(), "TT", 1.99, 10, models.Buy)
	sellOrder := models.NewOrder(uuid.NewV4(), "TT", 2.99, 5, models.Sell)

	ap := trading.NewOrderAppender()
	ap.Append(book, buyOrder)
	ap.Append(book, sellOrder)

	handler := NewOrderBookHandler(book)

	request, _ := http.NewRequest(http.MethodGet, "/orderbook/TT", nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodGet, "/orderbook/:symbol", handler.GetSymbolHandler)

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusOK)
	require.Equal(response.Header().Get("Content-Type"), "application/json")

	var symbol SymbolResponse

	body, err := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(body, &symbol)
	require.Nil(err)

	require.Equal("TT", symbol.Symbol, string(body))
	require.Len(symbol.Prices, 2)
	require.Equal(1.99, symbol.Prices[0].Price)
	require.Equal(uint(10), symbol.Prices[0].BuyQuantity)
	require.Equal(uint(1), symbol.Prices[0].BuyDepth)
	require.Equal(uint(0), symbol.Prices[0].SellQuantity)
	require.Equal(uint(0), symbol.Prices[0].SellDepth)

	require.Equal(2.99, symbol.Prices[1].Price)
	require.Equal(uint(0), symbol.Prices[1].BuyQuantity)
	require.Equal(uint(0), symbol.Prices[1].BuyDepth)
	require.Equal(uint(5), symbol.Prices[1].SellQuantity)
	require.Equal(uint(1), symbol.Prices[1].SellDepth)
}

func TestGetSymbolsHandler(t *testing.T) {

	require := require.New(t)
	book := models.NewOrderBook()

	buyOrderTT := models.NewOrder(uuid.NewV4(), "TT", 1.99, 10, models.Buy)
	sellOrderTT := models.NewOrder(uuid.NewV4(), "TT", 2.99, 5, models.Sell)
	buyOrderETE := models.NewOrder(uuid.NewV4(), "ETE", 1.99, 10, models.Buy)
	sellOrderETE := models.NewOrder(uuid.NewV4(), "ETE", 2.99, 5, models.Sell)

	ap := trading.NewOrderAppender()
	ap.Append(book, buyOrderTT)
	ap.Append(book, sellOrderTT)
	ap.Append(book, buyOrderETE)
	ap.Append(book, sellOrderETE)

	handler := NewOrderBookHandler(book)

	request, _ := http.NewRequest(http.MethodGet, "/orderbook", nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodGet, "/orderbook", handler.GetSymbolsHandler)

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusOK)
	require.Equal(response.Header().Get("Content-Type"), "application/json")

	var symbols []SymbolResponse

	body, err := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(body, &symbols)
	require.Nil(err)

	require.Len(symbols, 2)

	require.Len(symbols[0].Prices, 2)
	require.Len(symbols[1].Prices, 2)
}
