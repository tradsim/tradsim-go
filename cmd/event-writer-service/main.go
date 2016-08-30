package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/mantzas/incata"
	"github.com/mantzas/incata/marshal"
	incatamodel "github.com/mantzas/incata/model"
	"github.com/mantzas/incata/storage"
	"github.com/mantzas/incata/writer"
	uuid "github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/events"
)

func main() {
	//TODO: Configuration handling
	var url = "amqp://guest:guest@localhost:5672/tradsim"
	var subQueue = "order_events"
	var subExchange = "order_events"
	var pubExchange = "order_event_stored"
	var connectionString = "postgres://postgres:1234@localhost/orderevents?sslmode=disable"
	var dbName = "orderevents"

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC | log.Lshortfile)
	log.SetPrefix("ews ")

	setupIncata(connectionString, dbName)

	processor := events.NewOrderEventProcessor(url, subExchange, subQueue, pubExchange, processEnvelope)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		processor.Close()
		log.Printf("Event writer service stopped.")
		os.Exit(1)
	}()

	err := processor.Open()
	if err != nil {
		log.Fatalf("Failed to open processor! %s", err)
	}

	err = processor.Process()
	if err != nil {
		log.Fatalf("Processor failed to process! %s", err)
	}

	log.Printf("Event writer service exiting")
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

func processEnvelope(envelope *events.OrderEventEnvelope) (string, error) {

	untypedEvent, err := envelope.GetOrderEvent()
	if err != nil {
		return "", err
	}

	var sourceID uuid.UUID
	var occured time.Time
	var version uint

	switch envelope.EventType {
	case events.OrderAcceptedType:
		event := untypedEvent.(events.OrderAccepted)
		sourceID, occured, version, err = getOrderEventData(event.OrderEvent)
		log.Printf("Order accepted received: %s", event.String())
	case events.OrderAmendedType:
		event := untypedEvent.(events.OrderAmended)
		sourceID, occured, version, err = getOrderEventData(event.OrderEvent)
		log.Printf("Order amended received: %s", event.String())
	case events.OrderCancelledType:
		event := untypedEvent.(events.OrderCancelled)
		sourceID, occured, version, err = getOrderEventData(event.OrderEvent)
		log.Printf("Order cancelled received: %s", event.String())
	case events.OrderTradedType:
		event := untypedEvent.(events.OrderTraded)
		sourceID, occured, version, err = getOrderEventData(event.OrderEvent)
		log.Printf("Order traded received: %s", event.String())
	default:
		return "", errors.New("invalid order event type received")
	}

	if err != nil {
		return "", err
	}

	dbEvent := incatamodel.NewEvent(sourceID, occured, envelope.Payload, string(envelope.EventType), int(version))

	appender, err := incata.NewAppender()

	if err != nil {
		log.Printf("Faile to create a appender! %s", err)
		return "", err
	}

	return sourceID.String(), appender.Append(*dbEvent)
}

func getOrderEventData(orderEvent events.OrderEvent) (uuid.UUID, time.Time, uint, error) {
	sourceID, err := uuid.FromString(orderEvent.OrderID)
	if err != nil {
		return uuid.Nil, time.Now(), 0, err
	}

	return sourceID, orderEvent.Occured, orderEvent.Version, nil
}
