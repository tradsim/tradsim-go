package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"fmt"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/adaptlog"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tradsim/tradsim-go/cmd/exchange-service/trading"
	"github.com/tradsim/tradsim-go/mocks"
	"github.com/tradsim/tradsim-go/models"
)

func TestMain(m *testing.M) {
	adaptlog.ConfigureStdLevelLogger(adaptlog.DebugLevel, nil, "")
	retCode := m.Run()
	os.Exit(retCode)
}

func TestOrderCancelHandleAccepted(t *testing.T) {

	require := require.New(t)
	book := models.NewOrderBook()

	order := models.NewOrder(uuid.NewV4(), "TT", 1.99, 10, models.Buy)

	ap := trading.NewOrderAppender()
	ap.Append(book, order)

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockTrader{}, &mocks.MockCanceller{true}, &mocks.MockPublisher{})

	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/orders/%s", order.ID), nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodDelete, "/orders/:orderid", handler.OrderCancelHandle)

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusAccepted)
}

func TestOrderCancelHandleNotFound(t *testing.T) {

	require := require.New(t)
	book := models.NewOrderBook()

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockTrader{}, &mocks.MockCanceller{false}, &mocks.MockPublisher{})

	request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/orders/%s", uuid.NewV4()), nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodDelete, "/orders/:orderid", handler.OrderCancelHandle)

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusNotFound)
}

func TestOrderCancelHandleBadRequest(t *testing.T) {

	require := require.New(t)
	book := models.NewOrderBook()

	handler := NewOrderHandler(book, &mocks.MockAppender{}, &mocks.MockTrader{}, &mocks.MockCanceller{false}, &mocks.MockPublisher{})

	request, _ := http.NewRequest(http.MethodDelete, "/orders/123", nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodDelete, "/orders/:orderid", handler.OrderCancelHandle)

	router.ServeHTTP(response, request)

	require.Equal(response.Code, http.StatusBadRequest)
}
