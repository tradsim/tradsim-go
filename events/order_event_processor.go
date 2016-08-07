package events

import (
	"encoding/json"

	"github.com/mantzas/adaptlog"
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
	url        string
	exchange   string
	queue      string
	connection *amqp.Connection
	channel    *amqp.Channel
	logger     adaptlog.LevelLogger
	processor  func(envelope *OrderEventEnvelope) error
}

// NewOrderEventProcessor creates a new order event processor
func NewOrderEventProcessor(url string, exchange string, queue string, processor func(envelope *OrderEventEnvelope) error) *OrderEventProcessor {
	return &OrderEventProcessor{url, exchange, queue, nil, nil, adaptlog.NewStdLevelLogger("OrderEventProcessor"), processor}
}

// Open handles the opening of connection, channel, echange and queue
func (p *OrderEventProcessor) Open() error {

	conn, err := p.setupConnection(p.url)
	if err != nil {
		return err
	}

	p.connection = conn

	ch, err := p.setupExchangeAndQueue(conn, p.exchange, p.queue)
	if err != nil {
		return err
	}
	p.channel = ch

	return nil
}

// Close handles the closing of connection and channel
func (p *OrderEventProcessor) Close() {

	if p.connection != nil {
		p.connection.Close()
	}

	if p.channel != nil {
		p.channel.Close()
	}
}

// Process starts processing events
func (p *OrderEventProcessor) Process() error {

	msgs, err := p.getSubscriptionChannel(p.channel, p.queue)

	if err != nil {
		return err
	}

	for d := range msgs {

		var envelope OrderEventEnvelope

		err := json.Unmarshal(d.Body, &envelope)
		if err != nil {
			return err
		}

		err = p.processor(&envelope)
		if err == nil {
			p.logger.Errorf("Failed to process envelope %s", err)
		} else {
			d.Ack(false)
		}
	}

	return nil
}

func (p *OrderEventProcessor) setupConnection(url string) (*amqp.Connection, error) {

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	p.logger.Info("ampq: connection setup")
	return conn, nil
}

func (p *OrderEventProcessor) setupExchangeAndQueue(conn *amqp.Connection, exchange string, queue string) (*amqp.Channel, error) {

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

	p.logger.Infof("ampq: exchange %s declared", exchange)

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

	p.logger.Infof("ampq: queue %s declared", queue)

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil)
	if err != nil {
		return nil, err
	}

	p.logger.Infof("ampq: exchange %s with queue %s bound", exchange, queue)

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

	p.logger.Infof("ampq: subscription on queue %s set", queue)

	return msgs, nil
}
