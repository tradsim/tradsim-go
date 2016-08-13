package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/adaptlog"
	uuid "github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/cmd/exchange-service/trading"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

// OrderCreate model
type OrderCreate struct {
	ID        string  `json:"id"`        // unique id (uuid)
	Symbol    string  `json:"symbol"`    // symbol
	Quantity  uint    `json:"quantity"`  // quantity
	Direction string  `json:"direction"` // buy or sell
	Price     float64 `json:"price"`     // price
}

// OrderHandler handles orders
type OrderHandler struct {
	book      *models.OrderBook
	appender  trading.Appender
	trader    trading.Trader
	canceller trading.Canceller
	publisher events.EventPublisher
	logger    adaptlog.LevelLogger
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(book *models.OrderBook, appender trading.Appender, trader trading.Trader, canceller trading.Canceller, publisher events.EventPublisher) *OrderHandler {
	return &OrderHandler{book, appender, trader, canceller, publisher, adaptlog.NewStdLevelLogger("OrderHandler")}
}

// OrderCreateHandle is the handler for the orders
func (oh *OrderHandler) OrderCreateHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var orderCreate OrderCreate

	err := json.NewDecoder(r.Body).Decode(&orderCreate)

	if err != nil {
		oh.logger.Errorf("Failed to bind model! %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	direction, err := models.TradeDirectionFromString(orderCreate.Direction)
	if err != nil {
		oh.logger.Errorf("Failed to getting trade direction! %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orderID, err := uuid.FromString(orderCreate.ID)
	if err != nil {
		oh.logger.Errorf("Failed to getting order id! %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	order := models.NewOrder(orderID, orderCreate.Symbol, orderCreate.Price, orderCreate.Quantity, direction)

	createdEvent := events.NewOrderCreated(order.ID.String(), time.Now().UTC(), order.Symbol, order.Price, order.Quantity, order.Direction, 1)
	envelope, err := events.NewOrderEventEnvelope(createdEvent, createdEvent.EventType)
	if err != nil {
		oh.logger.Errorf("Failed to create order created event envelope! %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	oh.publisher.Publish(envelope)

	oh.trader.Trade(oh.book, order)

	if order.Status.IsTradeable() {
		err = oh.appender.Append(oh.book, order)
		if err != nil {
			oh.logger.Errorf("Failed to append order! %s", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

// OrderCancelHandle is the handler for cancelling a order
func (oh *OrderHandler) OrderCancelHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	orderID := uuid.FromStringOrNil(strings.ToUpper(p.ByName("orderid")))

	if orderID == uuid.Nil {
		oh.logger.Error("Failed to get orderID!")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	cancelled := oh.canceller.Cancel(oh.book, orderID)

	if cancelled {
		w.WriteHeader(http.StatusAccepted)
		oh.logger.Debugf("OrderCancelHandle: Order %s cancelled", orderID.String())
	} else {
		http.NotFound(w, r)
		oh.logger.Debugf("OrderCancelHandle: Order %s not found", orderID.String())
	}

}
