package events

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

// EventPublisher interface
type EventPublisher interface {
	Open() error
	Close()
	Publish(envelop *OrderEventEnvelope) error
}

// RabbitMqEventPublisher create implementation of a rabbit mq order event publisher
type RabbitMqEventPublisher struct {
	channel    *amqp.Channel
	connection *amqp.Connection
	exchange   string
	url        string
}

// NewRabbitMqEventPublisher creates a new rabbit mq event publisher
func NewRabbitMqEventPublisher(url string, exchange string) *RabbitMqEventPublisher {

	return &RabbitMqEventPublisher{exchange: exchange, url: url}
}

// Open handles the opening of a connection, setting up a channel and declaring a exchange
func (p *RabbitMqEventPublisher) Open() error {

	conn, err := amqp.Dial(p.url)
	if err != nil {
		return err
	}
	p.connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	p.channel = ch

	err = ch.ExchangeDeclare(
		p.exchange, // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)

	if err != nil {
		return err
	}

	return nil
}

// Close handles the connection and channel closing
func (p *RabbitMqEventPublisher) Close() {

	if p.channel != nil {
		p.channel.Close()
	}

	if p.connection != nil {
		p.connection.Close()
	}
}

// Publish publishes a order event envelope
func (p *RabbitMqEventPublisher) Publish(envelope *OrderEventEnvelope) error {

	jsonBytes, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		p.exchange, // exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBytes,
		})

	if err != nil {
		return err
	}

	return nil
}
