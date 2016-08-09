package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/mantzas/adaptlog"
	"github.com/tradsim/tradsim-go/cmd/exchange-service/handlers"
	"github.com/tradsim/tradsim-go/cmd/exchange-service/trading"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
	common_http "github.com/tradsim/tradsim-go/net/http"
)

func main() {

	adaptlog.ConfigureStdLevelLogger(adaptlog.DebugLevel, nil, "main")
	var url = "amqp://guest:guest@localhost:5672/tradsim"
	var exchange = "order_events"
	publisher := events.NewRabbitMqEventPublisher(url, exchange)

	err := publisher.Open()
	if err != nil {
		adaptlog.Level.Errorf("Failed to open publisher connection! %s", err)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		adaptlog.Level.Infoln("Exchange service stopped.")
		os.Exit(1)
	}()

	orderBook := models.NewOrderBook()
	appender := trading.NewOrderAppender()
	trader := trading.NewOrderTrader(publisher)
	orderHandler := handlers.NewOrderHandler(orderBook, appender, trader, publisher)
	// orderBookHandler := handlers.NewOrderBookHandler(orderBook)

	router := httprouter.New()

	router.POST("/orders", common_http.POSTJSONValidationMiddleware(orderHandler.OrderCreateHandle))
	// router.GET("/orderbook", common_http.GETValidationMiddleware(orderBookHandler.GetSymbolsHandler))
	// router.GET("/orderbook/:symbol", common_http.GETValidationMiddleware(orderBookHandler.GetSymbolHandler))

	adaptlog.Level.Info("Starting exchange  service.")

	adaptlog.Level.Fatal(http.ListenAndServe(":8081", router))
}
