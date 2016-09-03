package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/mantzas/incata"
	"github.com/mantzas/incata/marshal"
	"github.com/mantzas/incata/storage"
	"github.com/mantzas/incata/writer"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"github.com/tradsim/tradsim-go/events"
)

func main() {
	//TODO: Configuration handling
	var url = "amqp://guest:guest@localhost:5672/tradsim"
	var qn = "order_event_stored"
	var exc = "order_event_stored"
	var connectionString = "postgres://postgres:1234@localhost/orderevents?sslmode=disable"
	var dbName = "orderevents"

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC | log.Lshortfile)
	log.SetPrefix("eas ")

	setupIncata(connectionString, dbName)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("Event aggregation service stopped.")
		os.Exit(1)
	}()

	msgs, err := createSubscription(url, exc, qn)
	if err != nil {
		log.Fatalf("Failed to create subscription! %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			processDelivery(&d)
		}
	}()

	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever
	log.Fatal("Event aggregation service stopped unexpectedly.")
}

func processDelivery(d *amqp.Delivery) {
	var env events.OrderEventEnvelope

	err := json.Unmarshal(d.Body, &env)
	if err != nil {
		log.Fatalf("Failed to unmarshal envelope. %s", err)
	}

	untypedEvent, err := env.GetOrderEvent()
	if err != nil {
		log.Fatalf("Failed to get event. %s", err)
	}

	var event events.OrderEventStored

	switch env.EventType {
	case events.OrderEventStoredType:
		event = untypedEvent.(events.OrderEventStored)
	default:
		log.Fatalf("Invalid event received %s", env.EventType)
	}

	log.Printf("Received %s", event.EventType)
	err = processStoredEvent(event)
	if err != nil {
		log.Fatalf("Failed to process stored event. %s", err)
	}
	d.Ack(false)
}

func processStoredEvent(event events.OrderEventStored) error {

	sourceID, err := uuid.FromString(event.OrderID)
	if err != nil {
		return err
	}

	r, err := incata.NewRetriever()
	if err != nil {
		return err
	}

	_, err = r.Retrieve(sourceID)
	if err != nil {
		return err
	}

	// TODO: Aggregate Events

	// TODO: Store order and position to db

	return nil
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

func createSubscription(url string, exc string, qn string) (<-chan amqp.Delivery, error) {

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		qn,    // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		exc,    // exchange
		false,
		nil)

	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		return nil, err
	}

	return msgs, nil
}
