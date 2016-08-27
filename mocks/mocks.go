package mocks

import (
	"github.com/satori/go.uuid"
	"github.com/tradsim/tradsim-go/events"
	"github.com/tradsim/tradsim-go/models"
)

// MockPublisher mocks a publisher
type MockPublisher struct {
	Envelopes []*events.OrderEventEnvelope
}

// Open the connection
func (mp *MockPublisher) Open() error {
	return nil
}

// Close the connection
func (mp *MockPublisher) Close() {
}

// Publish a envelope
func (mp *MockPublisher) Publish(envelope *events.OrderEventEnvelope) error {
	mp.Envelopes = append(mp.Envelopes, envelope)
	return nil
}

// MockAppender for mocking the appender
type MockAppender struct {
}

// Append order
func (ma *MockAppender) Append(book *models.OrderBook, order *models.Order) error {
	return nil
}

// MockTrader for mocking the trader
type MockTrader struct {
}

// Trade order
func (mt *MockTrader) Trade(book *models.OrderBook, order *models.Order) {

}

// MockCanceller for mocking the canceller
type MockCanceller struct {
	Cancelled bool
}

// Cancel the order
func (mc *MockCanceller) Cancel(book *models.OrderBook, orderID uuid.UUID) bool {
	return mc.Cancelled
}

// MockAmender for mocking the amender
type MockAmender struct {
	Amended bool
}

// Amend the order
func (ma *MockAmender) Amend(book *models.OrderBook, order *models.Order) bool {
	return ma.Amended
}
