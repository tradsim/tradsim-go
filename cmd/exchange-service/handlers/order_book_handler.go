package handlers

import (
	"net/http"

	"encoding/json"

	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/adaptlog"
	"github.com/tradsim/tradsim-go/models"
)

// SymbolPriceResponse returns a price and thq buy sell quantities
type SymbolPriceResponse struct {
	Price        float64 `json:"price"`
	BuyQuantity  uint    `json:"buy_quantity"`
	BuyDepth     uint    `json:"buy_depth"`
	SellQuantity uint    `json:"sell_quantity"`
	SellDepth    uint    `json:"sell_depth"`
}

// SymbolResponse returns a symbol along with the prices and quantities
type SymbolResponse struct {
	Symbol string                `json:"symbol"`
	Prices []SymbolPriceResponse `json:"prices"`
}

// OrderBookHandler handles book requests
type OrderBookHandler struct {
	book   *models.OrderBook
	logger adaptlog.LevelLogger
}

// NewOrderBookHandler creates a new order book handler
func NewOrderBookHandler(book *models.OrderBook) *OrderBookHandler {
	return &OrderBookHandler{book, adaptlog.NewStdLevelLogger("OrderBookHandler")}
}

// GetSymbolsHandler is the handler for getting symbols from the order book
func (obh *OrderBookHandler) GetSymbolsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	obh.logger.Debug("GetSymbolsHandler: Request symbols")

	var symbols []SymbolResponse

	for key, prices := range obh.book.Symbols {

		symbols = append(symbols, getSymbolResponse(key, prices))
	}

	if len(symbols) == 0 {
		obh.logger.Debug("GetSymbolsHandler: Symbols not found")
		http.NotFound(w, r)
		return
	}

	encoded, _ := json.Marshal(symbols)
	w.Header().Set("Content-Type", "application/json")
	w.Write(encoded)
	obh.logger.Debugf("GetSymbolsHandler: Returned %d symbols", len(symbols))
}

// GetSymbolHandler is the handler for getting a symbol from the order book
func (obh *OrderBookHandler) GetSymbolHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	symbol := strings.ToUpper(p.ByName("symbol"))

	obh.logger.Debugf("GetSymbolHandler: Request %s", symbol)

	prices, ok := obh.book.Symbols[symbol]

	if !ok {
		obh.logger.Debugf("GetSymbolHandler: Symbol %s not found", symbol)
		http.NotFound(w, r)
		return
	}

	response := getSymbolResponse(symbol, prices)

	encoded, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(encoded)
	obh.logger.Debugf("GetSymbolHandler: Symbol %s returned", symbol)
}

func getSymbolResponse(symbol string, prices []*models.OrderPrice) SymbolResponse {

	response := SymbolResponse{symbol, make([]SymbolPriceResponse, 0)}

	for _, price := range prices {
		response.Prices = append(response.Prices, getSymbolPriceResponse(price))
	}

	return response
}

func getSymbolPriceResponse(price *models.OrderPrice) SymbolPriceResponse {
	return SymbolPriceResponse{price.Price, price.Buy.Quantity, uint(len(price.Buy.Orders)),
		price.Sell.Quantity, uint(len(price.Sell.Orders))}
}
