package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"fmt"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/adaptlog"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/cmd/exchange-service/trading"
	"github.com/tradsim/tradsim-go/mocks"
	"github.com/tradsim/tradsim-go/models"
	common_http "github.com/tradsim/tradsim-go/net/http"
)

func TestMain(m *testing.M) {
	adaptlog.ConfigureStdLevelLogger(adaptlog.DebugLevel, nil, "")
	retCode := m.Run()
	os.Exit(retCode)
}

func TestOrderCreateHandle(t *testing.T) {
	require := require.New(t)
	book := models.NewOrderBook()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 1.99, 10, models.Buy)
	createOrder := OrderDTO{uuid.NewV4().String(), "TT", 10, models.Sell.String(), 1.99}

	ap := trading.NewOrderAppender()
	ap.Append(book, order1)

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockAmender{}, &mocks.MockTrader{}, &mocks.MockCanceller{Cancelled: true}, &mocks.MockPublisher{})

	encodedOrder, _ := json.Marshal(createOrder)

	request, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(encodedOrder))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodPost, "/orders", common_http.DefaultPOSTJSONValidationMiddleware(handler.OrderCreateHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusAccepted, response.Code)
}

func TestOrderCreateHandleInvalidPayloadBadRequest(t *testing.T) {
	require := require.New(t)
	book := models.NewOrderBook()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 1.99, 10, models.Buy)

	ap := trading.NewOrderAppender()
	ap.Append(book, order1)

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockAmender{}, &mocks.MockTrader{}, &mocks.MockCanceller{Cancelled: true}, &mocks.MockPublisher{})

	encodedOrder, _ := json.Marshal("{\"test\":123}")

	request, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(encodedOrder))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodPost, "/orders", common_http.DefaultPOSTJSONValidationMiddleware(handler.OrderCreateHandle))

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusBadRequest)
}

func TestOrderCreateHandleInvalidOrderIDBadRequest(t *testing.T) {
	require := require.New(t)
	book := models.NewOrderBook()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 1.99, 10, models.Buy)

	ap := trading.NewOrderAppender()
	ap.Append(book, order1)

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockAmender{}, &mocks.MockTrader{}, &mocks.MockCanceller{Cancelled: true}, &mocks.MockPublisher{})

	createOrder := OrderDTO{"XXX", "TT", 10, models.Sell.String(), 1.99}
	encodedOrder, _ := json.Marshal(createOrder)

	request, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(encodedOrder))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodPost, "/orders", common_http.DefaultPOSTJSONValidationMiddleware(handler.OrderCreateHandle))

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusBadRequest)
}

func TestOrderCreateHandleInvalidTradeDirectionBadRequest(t *testing.T) {
	require := require.New(t)
	book := models.NewOrderBook()

	order1 := models.NewOrder(uuid.NewV4(), "TT", 1.99, 10, models.Buy)
	createOrder := OrderDTO{uuid.NewV4().String(), "TT", 10, "XXX", 1.99}

	ap := trading.NewOrderAppender()
	ap.Append(book, order1)

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockAmender{}, &mocks.MockTrader{}, &mocks.MockCanceller{Cancelled: true}, &mocks.MockPublisher{})

	encodedOrder, _ := json.Marshal(createOrder)

	request, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(encodedOrder))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodPost, "/orders", common_http.DefaultPOSTJSONValidationMiddleware(handler.OrderCreateHandle))

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusBadRequest)
}

func TestOrderCancelHandleAccepted(t *testing.T) {

	require := require.New(t)
	book := models.NewOrderBook()

	order := models.NewOrder(uuid.NewV4(), "TT", 1.99, 10, models.Buy)

	ap := trading.NewOrderAppender()
	ap.Append(book, order)

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockAmender{}, &mocks.MockTrader{}, &mocks.MockCanceller{Cancelled: true}, &mocks.MockPublisher{})

	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/orders/%s", order.ID), nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodDelete, "/orders/:orderid", common_http.DefaultDELETEValidationMiddleware(handler.OrderCancelHandle))

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusAccepted)
}

func TestOrderCancelHandleNotFound(t *testing.T) {

	require := require.New(t)
	book := models.NewOrderBook()

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockAmender{}, &mocks.MockTrader{}, &mocks.MockCanceller{Cancelled: false}, &mocks.MockPublisher{})

	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/orders/%s", uuid.NewV4()), nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodDelete, "/orders/:orderid", common_http.DefaultDELETEValidationMiddleware(handler.OrderCancelHandle))

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusNotFound)
}

func TestOrderCancelHandleBadRequest(t *testing.T) {

	require := require.New(t)
	book := models.NewOrderBook()

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockAmender{}, &mocks.MockTrader{}, &mocks.MockCanceller{Cancelled: false}, &mocks.MockPublisher{})

	request, _ := http.NewRequest(http.MethodDelete, "/orders/123", nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodDelete, "/orders/:orderid", common_http.DefaultDELETEValidationMiddleware(handler.OrderCancelHandle))

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusBadRequest)
}

func TestOrderAmendHandleAccepted(t *testing.T) {
	require := require.New(t)
	book := models.NewOrderBook()

	amendOrder := OrderDTO{uuid.NewV4().String(), "TT", 10, models.Sell.String(), 1.99}

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockAmender{Amended: true}, &mocks.MockTrader{}, &mocks.MockCanceller{}, &mocks.MockPublisher{})

	encodedOrder, _ := json.Marshal(amendOrder)

	request, _ := http.NewRequest(http.MethodPut, "/orders", bytes.NewBuffer(encodedOrder))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodPut, "/orders", common_http.DefaultPUTJSONValidationMiddleware(handler.OrderAmendHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusAccepted, response.Code)
}

func TestOrderAmendHandleBadRequest(t *testing.T) {
	require := require.New(t)
	book := models.NewOrderBook()

	amendOrder := OrderDTO{uuid.NewV4().String(), "TT", 10, models.Sell.String(), 1.99}

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockAmender{Amended: false}, &mocks.MockTrader{}, &mocks.MockCanceller{}, &mocks.MockPublisher{})

	encodedOrder, _ := json.Marshal(amendOrder)

	request, _ := http.NewRequest(http.MethodPut, "/orders", bytes.NewBuffer(encodedOrder))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodPut, "/orders", common_http.DefaultPUTJSONValidationMiddleware(handler.OrderAmendHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusBadRequest, response.Code)
}
