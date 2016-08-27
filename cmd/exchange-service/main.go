package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/tradsim/tradsim-go/cmd/exchange-service/handlers"
	"github.com/tradsim/tradsim-go/cmd/exchange-service/trading"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
	common_http "github.com/tradsim/tradsim-go/net/http"
)

func main() {

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC | log.Lshortfile)
	log.SetPrefix("es ")

	var url = "amqp://guest:guest@localhost:5672/tradsim"
	var exchange = "order_events"
	publisher := events.NewRabbitMqEventPublisher(url, exchange)

	err := publisher.Open()
	if err != nil {
		log.Printf("Failed to open publisher connection! %s", err)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("Exchange service stopped.")
		os.Exit(1)
	}()

	orderBook := models.NewOrderBook()
	appender := trading.NewOrderAppender()
	amender := trading.NewOrderAmender(publisher)
	trader := trading.NewOrderTrader(publisher)
	canceller := trading.NewOrderCanceller(publisher)
	orderHandler := handlers.NewOrderHandler(orderBook, appender, amender, trader, canceller, publisher)
	orderBookHandler := handlers.NewOrderBookHandler(orderBook)

	router := httprouter.New()

	router.POST("/orders", common_http.POSTJSONValidationMiddleware(orderHandler.OrderCreateHandle))
	router.PUT("/orders", common_http.PUTJSONValidationMiddleware(orderHandler.OrderAmendHandle))
	router.DELETE("/orders/:orderid", common_http.DELETEValidationMiddleware(orderHandler.OrderCancelHandle))
	router.GET("/orderbook", common_http.GETValidationMiddleware(orderBookHandler.GetSymbolsHandler))
	router.GET("/orderbook/:symbol", common_http.GETValidationMiddleware(orderBookHandler.GetSymbolHandler))

	log.Print("Starting exchange service.")

	log.Fatal(http.ListenAndServe(":8081", router))
}
