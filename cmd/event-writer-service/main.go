package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/mantzas/adaptlog"
	"github.com/mantzas/incata"
	"github.com/mantzas/incata/marshal"
	incatamodel "github.com/mantzas/incata/model"
	"github.com/mantzas/incata/storage"
	"github.com/mantzas/incata/writer"
	"github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/events"
)

func main() {
	//TODO: Configuration handling
	var url = "amqp://guest:guest@localhost:5672/tradsim"
	var queue = "order_events"
	var exchange = "order_events"
	var connectionString = "postgres://postgres:1234@localhost/orderevents?sslmode=disable"
	var dbName = "orderevents"

	adaptlog.ConfigureStdLevelLogger(adaptlog.DebugLevel, nil, "main")

	setupIncata(connectionString, dbName)

	processor := events.NewOrderEventProcessor(url, exchange, queue, processEnvelope)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		processor.Close()
		adaptlog.Level.Infoln("Event writer service stopped.")
		os.Exit(1)
	}()

	err := processor.Open()
	if err != nil {
		adaptlog.Level.Errorf("Failed to open processor! %s", err)
		return
	}

	err = processor.Process()
	if err != nil {
		adaptlog.Level.Errorf("Processor failed to process! %s", err)
		return
	}

	adaptlog.Level.Infoln("Event writer service exiting")
}

func setupIncata(connection string, dbName string) {

	storage, err := storage.NewStorage(storage.PostgreSQL, connection, dbName)

	if err != nil {
		panic(err)
	}

	sr := marshal.NewJSONMarshaller()
	wr := writer.NewSQLWriter(storage, sr)

	incata.SetupAppender(wr)
}

func processEnvelope(envelope *events.OrderEventEnvelope) error {

	untypedEvent, err := envelope.GetOrderEvent()
	if err != nil {
		return err
	}

	var sourceID uuid.UUID
	var created time.Time
	var version uint

	switch envelope.EventType {
	case events.OrderCreatedType:
		event := untypedEvent.(events.OrderCreated)
		sourceID, created, version, err = getOrderEventData(event.OrderEvent)
		adaptlog.Level.Infof("Order created received: %s", event.String())
	case events.OrderAmendedType:
		event := untypedEvent.(events.OrderAmended)
		sourceID, created, version, err = getOrderEventData(event.OrderEvent)
		adaptlog.Level.Infof("Order amended received: %s", event.String())
	case events.OrderCancelledType:
		event := untypedEvent.(events.OrderCancelled)
		sourceID, created, version, err = getOrderEventData(event.OrderEvent)
		adaptlog.Level.Infof("Order cancelled received: %s", event.String())
	case events.OrderTradedType:
		event := untypedEvent.(events.OrderTraded)
		sourceID, created, version, err = getOrderEventData(event.OrderEvent)
		adaptlog.Level.Infof("Order traded received: %s", event.String())
	default:
		return errors.New("invalid order event type received")
	}

	if err != nil {
		return err
	}

	dbEvent := incatamodel.NewEvent(sourceID, created, envelope.Payload, string(envelope.EventType), int(version))

	appender, err := incata.NewAppender()

	if err != nil {
		adaptlog.Level.Errorf("Faile to create a appender! %s", err)
		return err
	}

	return appender.Append(*dbEvent)
}

func getOrderEventData(orderEvent events.OrderEvent) (uuid.UUID, time.Time, uint, error) {
	sourceID, err := uuid.FromString(orderEvent.OrderID)
	if err != nil {
		return uuid.Nil, time.Now(), 0, err
	}

	return sourceID, orderEvent.Occured, orderEvent.Version, nil
}
