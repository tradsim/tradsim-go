package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/mantzas/incata"
	"github.com/mantzas/incata/marshal"
	"github.com/mantzas/incata/reader"
	"github.com/mantzas/incata/storage"
	"github.com/streadway/amqp"
	"github.com/tradsim/tradsim-go/cmd/event-aggregation-service/aggregator"
	"github.com/tradsim/tradsim-go/cmd/event-aggregation-service/data"
	"github.com/tradsim/tradsim-go/cmd/event-aggregation-service/processor"
	"github.com/tradsim/tradsim-go/events"
)

func main() {
	//TODO: Configuration handling
	var url = "amqp://guest:guest@localhost:5672/tradsim"
	var qn = "order_event_stored"
	var exc = "order_event_stored"
	var evCon = "postgres://postgres:1234@localhost/orderevents?sslmode=disable"
	var dbName = "orderevents"
	var orCon = "postgres://postgres:1234@localhost/order?sslmode=disable"

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC | log.Lshortfile)
	log.SetPrefix("eas ")

	setupIncata(evCon, dbName)

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
	rt, err := incata.NewRetriever()
	if err != nil {
		log.Fatalf("Failed to create event retriever! %s", err)
	}

	orDB, err := getOrderDb(orCon)
	if err != nil {
		log.Fatalf("Failed to connect to order db! %s", err)
	}

	repo, err := data.NewOrderRepository(orDB)
	if err != nil {
		log.Fatalf("Failed to create order repo! %s", err)
	}

	evagg := aggregator.NewEventAggregator()
	oragg := aggregator.NewOrderAggregator()
	prc := processor.NewEventProcessor(rt, evagg, oragg, repo)

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			go func(prc processor.Processor, msg *amqp.Delivery) {
				processDelivery(prc, msg)
			}(prc, &msg)
		}
	}()

	log.Printf(" [*] Waiting for events. To exit press CTRL+C")
	<-forever
	log.Fatal("Event aggregation service stopped unexpectedly.")
}

func processDelivery(prc processor.Processor, d *amqp.Delivery) {
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

	err = prc.Process(event)
	if err != nil {
		log.Fatalf("Failed to process stored event. %s", err)
	}
	d.Ack(false)
}

func setupIncata(connection string, dbName string) {

	storage, err := storage.NewStorage(storage.PostgreSQL, connection, dbName)

	if err != nil {
		panic(err)
	}

	sr := marshal.NewJSONMarshaller()
	rd := reader.NewSQLReader(storage, sr)

	incata.SetupRetriever(rd)
}

func getOrderDb(cn string) (*sql.DB, error) {

	db, err := sql.Open("postgres", cn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
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
