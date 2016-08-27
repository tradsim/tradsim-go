package events

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// EventProcessor interface
type EventProcessor interface {
	Open() error
	Close()
	Process() error
}

// OrderEventProcessor defines the order event struct
type OrderEventProcessor struct {
	url       string
	subExc    string
	subQ      string
	pubExc    string
	conn      *amqp.Connection
	subCh     *amqp.Channel
	pubCh     *amqp.Channel
	processor func(envelope *OrderEventEnvelope) (string, error)
}

// NewOrderEventProcessor creates a new order event processor
func NewOrderEventProcessor(url string, subExchange string, subQueue string, pubExchange string, processor func(envelope *OrderEventEnvelope) (string, error)) *OrderEventProcessor {
	return &OrderEventProcessor{url, subExchange, subQueue, pubExchange, nil, nil, nil, processor}
}

// Open handles the opening of connection, channel, echange and queue
func (p *OrderEventProcessor) Open() error {

	conn, err := p.setupConnection(p.url)
	if err != nil {
		return err
	}

	p.conn = conn

	subCh, err := p.setupSubscribeExchangeAndQueue(conn, p.subExc, p.subQ)
	if err != nil {
		return err
	}
	p.subCh = subCh

	pubCh, err := p.setupPublishChannel(p.pubExc)
	if err != nil {
		return err
	}
	p.pubCh = pubCh

	return nil
}

// Close handles the closing of connection and channel
func (p *OrderEventProcessor) Close() {

	if p.conn != nil {
		p.conn.Close()
	}

	if p.subCh != nil {
		p.subCh.Close()
	}
}

// Process starts processing events
func (p *OrderEventProcessor) Process() error {

	msgs, err := p.getSubscriptionChannel(p.subCh, p.subQ)

	if err != nil {
		return err
	}

	for d := range msgs {

		var envelope OrderEventEnvelope

		err := json.Unmarshal(d.Body, &envelope)
		if err != nil {
			return err
		}

		orderID, err := p.processor(&envelope)
		if err != nil {
			log.Printf("Failed to process envelope %s", err)
		} else {
			d.Ack(false)
			p.publishOrderEventStored(orderID)
		}
	}

	return nil
}

func (p *OrderEventProcessor) setupConnection(url string) (*amqp.Connection, error) {

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	log.Print("ampq: connection setup")
	return conn, nil
}

func (p *OrderEventProcessor) setupSubscribeExchangeAndQueue(conn *amqp.Connection, exchange string, queue string) (*amqp.Channel, error) {

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchange, // name
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

	log.Printf("ampq: exchange %s declared", exchange)

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return nil, err
	}

	log.Printf("ampq: queue %s declared", queue)

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil)
	if err != nil {
		return nil, err
	}

	log.Printf("ampq: exchange %s with queue %s bound", exchange, queue)
	return ch, nil
}

func (p *OrderEventProcessor) getSubscriptionChannel(ch *amqp.Channel, queue string) (<-chan amqp.Delivery, error) {

	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	log.Printf("ampq: subscription on queue %s set", queue)

	return msgs, nil
}

func (p *OrderEventProcessor) setupPublishChannel(pubExchange string) (*amqp.Channel, error) {
	ch, err := p.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		pubExchange, // name
		"fanout",    // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)

	return ch, nil
}

func (p *OrderEventProcessor) publishOrderEventStored(orderID string) error {
	event := NewOrderEventStored(orderID, time.Now().UTC(), 1)
	envelope, err := NewOrderEventEnvelope(event, event.EventType)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	err = p.pubCh.Publish(
		p.pubExc, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBytes,
		})

	if err != nil {
		log.Print("Failed to publish stored event")
	}

	return err
}
