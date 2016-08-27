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

// OrderDTO model
type OrderDTO struct {
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
	amender   trading.Amender
	trader    trading.Trader
	canceller trading.Canceller
	publisher events.EventPublisher
	logger    adaptlog.LevelLogger
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(book *models.OrderBook, appender trading.Appender, amender trading.Amender, trader trading.Trader, canceller trading.Canceller, publisher events.EventPublisher) *OrderHandler {
	return &OrderHandler{book, appender, amender, trader, canceller, publisher, adaptlog.NewStdLevelLogger("OrderHandler")}
}

// OrderCreateHandle is the handler for the orders
func (oh *OrderHandler) OrderCreateHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	dto, direction, orderID, err := oh.getPayloadData(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	order := models.NewOrder(orderID, dto.Symbol, dto.Price, dto.Quantity, direction)

	acceptedEvent := events.NewOrderAccepted(order.ID.String(), time.Now().UTC(), order.Symbol, order.Price, order.Quantity, order.Direction, 1)
	envelope, err := events.NewOrderEventEnvelope(acceptedEvent, acceptedEvent.EventType)
	if err != nil {
		oh.logger.Errorf("Failed to create order accepted event envelope! %s", err)
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

// OrderAmendHandle is the handler for amending a order
func (oh *OrderHandler) OrderAmendHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	dto, direction, orderID, err := oh.getPayloadData(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	order := models.NewOrder(orderID, dto.Symbol, dto.Price, dto.Quantity, direction)

	amended := oh.amender.Amend(oh.book, order)

	if amended {
		w.WriteHeader(http.StatusAccepted)
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
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

func (oh *OrderHandler) getPayloadData(r *http.Request) (OrderDTO, models.TradeDirection, uuid.UUID, error) {

	var dto OrderDTO
	var direction models.TradeDirection
	var orderID uuid.UUID

	err := json.NewDecoder(r.Body).Decode(&dto)

	if err != nil {
		oh.logger.Errorf("Failed to bind model! %s", err)
		return dto, direction, orderID, err
	}

	direction, err = models.TradeDirectionFromString(dto.Direction)
	if err != nil {
		oh.logger.Errorf("Failed to getting trade direction! %s", err)
		return dto, direction, orderID, err
	}

	orderID, err = uuid.FromString(dto.ID)
	if err != nil {
		oh.logger.Errorf("Failed to getting order id! %s", err)
		return dto, direction, orderID, err
	}

	return dto, direction, orderID, nil
}
