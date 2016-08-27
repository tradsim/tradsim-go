package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
)

func TestGETValidationMiddlewareSuccess(t *testing.T) {

	require := require.New(t)

	request, _ := http.NewRequest(http.MethodGet, "/orderbook/TT", nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodGet, "/orderbook/:symbol", GETValidationMiddleware(MockHandler{false}.MockHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusOK, response.Code)
}

func TestGETValidationMiddlewareFailure(t *testing.T) {

	require := require.New(t)

	request, _ := http.NewRequest(http.MethodPost, "/orderbook/TT", nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodPost, "/orderbook/:symbol", GETValidationMiddleware(MockHandler{false}.MockHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusBadRequest, response.Code)
}

func TestPOSTValidationMiddlewareSuccess(t *testing.T) {

	require := require.New(t)

	request, _ := http.NewRequest(http.MethodPost, "/orderbook/TT", nil)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodPost, "/orderbook/:symbol", POSTJSONValidationMiddleware(MockHandler{false}.MockHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusOK, response.Code)
}

func TestPOSTValidationMiddlewarWrongContentTypeFailure(t *testing.T) {

	require := require.New(t)

	request, _ := http.NewRequest(http.MethodPost, "/orderbook/TT", nil)

	request.Header.Set("Content-Type", "application/xml")

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodPost, "/orderbook/:symbol", POSTJSONValidationMiddleware(MockHandler{false}.MockHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusBadRequest, response.Code)
}

func TestPostValidationMiddlewareFailure(t *testing.T) {

	require := require.New(t)

	request, _ := http.NewRequest(http.MethodGet, "/orderbook/TT", nil)

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodGet, "/orderbook/:symbol", POSTJSONValidationMiddleware(MockHandler{false}.MockHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusBadRequest, response.Code)
}

func TestRecoveryMiddlewareSuccess(t *testing.T) {

	require := require.New(t)

	request, _ := http.NewRequest(http.MethodPost, "/orderbook/TT", nil)

	response := httptest.NewRecorder()

	router := httprouter.New()

	router.Handle(http.MethodPost, "/orderbook/:symbol", RecoveryMiddleware(MockHandler{true}.MockHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusInternalServerError, response.Code)
}

func TestLoggingMiddlewareSuccess(t *testing.T) {

	require := require.New(t)

	request, _ := http.NewRequest(http.MethodPost, "/orderbook/TT", nil)

	response := httptest.NewRecorder()

	router := httprouter.New()

	router.Handle(http.MethodPost, "/orderbook/:symbol", LoggingMiddleware(MockHandler{false}.MockHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusOK, response.Code)
}

func TestDefaultMiddlewareSuccess(t *testing.T) {

	require := require.New(t)

	request, _ := http.NewRequest(http.MethodPost, "/orderbook/TT", nil)

	response := httptest.NewRecorder()

	router := httprouter.New()

	router.Handle(http.MethodPost, "/orderbook/:symbol", DefaultMiddleware(MockHandler{false}.MockHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusOK, response.Code)
}

func TestDefaultPostJSONValidationMiddlewareSuccess(t *testing.T) {

	require := require.New(t)

	request, _ := http.NewRequest(http.MethodPost, "/orderbook/TT", nil)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router := httprouter.New()
	router.Handle(http.MethodPost, "/orderbook/:symbol", DefaultPOSTJSONValidationMiddleware(MockHandler{false}.MockHandle))

	router.ServeHTTP(response, request)

	require.Equal(http.StatusOK, response.Code)
}

type MockHandler struct {
	panic bool
}

func (m MockHandler) MockHandle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if m.panic {
		panic("TEST")
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte("test"))
}
