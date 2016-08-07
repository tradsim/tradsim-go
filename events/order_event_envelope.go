package events

import (
	"encoding/json"
	"errors"
)

// OrderEventType defines a order event type enumeration
type OrderEventType string

// Order Event type constants
const (
	OrderCreatedType  OrderEventType = "OrderCreated"
	OrderAmendedType  OrderEventType = "OrderAmended"
	OrderCanceledType OrderEventType = "OrderCanceled"
	OrderTradedType   OrderEventType = "OrderTraded"
)

// GetEventType returns the event type from a event
func GetEventType(event interface{}) (OrderEventType, error) {
	switch event.(type) {
	case OrderCreated:
		return OrderCreatedType, nil
	case OrderAmended:
		return OrderAmendedType, nil
	case OrderCanceled:
		return OrderCanceledType, nil
	case OrderTraded:
		return OrderTradedType, nil
	default:
		return "", errors.New("invalid event provided")
	}
}

// OrderEventEnvelope encapsulates the event send to a subscriber
type OrderEventEnvelope struct {
	EventType OrderEventType `json:"event_type"`
	Payload   string         `json:"payload"`
}

// NewOrderEventEnvelope creates a new order event envelope from a event
func NewOrderEventEnvelope(event interface{}, eventType OrderEventType) (*OrderEventEnvelope, error) {

	jsonBytes, err := json.Marshal(event)

	if err != nil {
		return nil, err
	}

	return &OrderEventEnvelope{eventType, string(jsonBytes)}, nil
}

// GetOrderEvent returns the specific order event
func (e *OrderEventEnvelope) GetOrderEvent() (interface{}, error) {

	switch e.EventType {
	case OrderCreatedType:
		return e.getCreatedEvent()
	case OrderAmendedType:
		return e.getAmendedEvent()
	case OrderCanceledType:
		return e.getCanceledEvent()
	case OrderTradedType:
		return e.getTradedEvent()
	default:
		return nil, errors.New("invalid order event type provided")
	}
}

func (e *OrderEventEnvelope) getCreatedEvent() (interface{}, error) {
	var event OrderCreated
	err := e.getEvent(e.Payload, &event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (e *OrderEventEnvelope) getAmendedEvent() (interface{}, error) {
	var event OrderAmended
	err := e.getEvent(e.Payload, &event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (e *OrderEventEnvelope) getCanceledEvent() (interface{}, error) {
	var event OrderCanceled
	err := e.getEvent(e.Payload, &event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (e *OrderEventEnvelope) getTradedEvent() (interface{}, error) {
	var event OrderTraded
	err := e.getEvent(e.Payload, &event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (e *OrderEventEnvelope) getEvent(payload string, event interface{}) error {

	err := json.Unmarshal([]byte(payload), &event)

	if err != nil {
		return err
	}

	return nil
}
